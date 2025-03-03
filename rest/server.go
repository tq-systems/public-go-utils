/*
 * REST utilities - server.go
 * Copyright (c) 2020 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Matthias Schiffer and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

package rest

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/tq-systems/public-go-utils/v3/auth"
	"github.com/tq-systems/public-go-utils/v3/log"
)

// Server structure
type Server struct {
	router   *mux.Router
	baseURL  string
	listener net.Listener
}

// Route represents an API endpoint
type Route struct {
	Method  string
	Pattern string
	Role    interface{}
	Handler func(r *http.Request) *Response
}

// Listener is the listener configuration structure
type Listener struct {
	Address string
	Proto   string
	Group   string
}

const (
	wsTimeout = 5 * time.Second
)

// MakeServer creates a server
func MakeServer(baseURL string) *Server {
	return &Server{
		router:  mux.NewRouter().UseEncodedPath(),
		baseURL: baseURL,
	}
}

// GetRouter returns the router
func (srv *Server) GetRouter() *mux.Router {
	return srv.router
}

// CheckAuth proofs authorization
func CheckAuth(role interface{}, authorization string) error {
	authSplit := strings.Split(authorization, " ")
	if len(authSplit) != 2 || authSplit[0] != "Bearer" {
		return errors.New("Invalid authorization token")
	}
	user, err := auth.ValidateAuthToken(authSplit[1])
	if err != nil {
		return err
	}

	if !user.HasRole(role) {
		return errors.New("Insufficient permissions")
	}

	return nil
}

func (srv *Server) handleAuthorized(role interface{}, handler func(w http.ResponseWriter, r *http.Request) *Response, w http.ResponseWriter, r *http.Request) *Response {
	if CheckAuth(role, r.Header.Get("Authorization")) != nil {
		return NewEmptyResponseWithStatus(http.StatusUnauthorized)
	}

	return handler(w, r)
}

// AddHandlerWithWriter adds a handler with writer
func (srv *Server) AddHandlerWithWriter(method string, pattern string, role interface{}, handler func(w http.ResponseWriter, r *http.Request) *Response) *mux.Route {
	return srv.AddAuthRouteWithWriter(method, pattern, role, handler)
}

// AddHandler is an alias for AddAuthRoute
func (srv *Server) AddHandler(method string, pattern string, role string, handler func(r *http.Request) *Response) *mux.Route {
	return srv.AddAuthRoute(method, pattern, role, handler)
}

// AddAuthRouteWithWriter adds authorization route with writer
func (srv *Server) AddAuthRouteWithWriter(method string, pattern string, role interface{}, handler func(w http.ResponseWriter, r *http.Request) *Response) *mux.Route {
	return srv.AddRouteWithWriter(method, pattern, func(w http.ResponseWriter, r *http.Request) *Response {
		return srv.handleAuthorized(role, handler, w, r)
	})
}

// AddAuthRoute adds a protected route
//
//	srv := rest.MakeServer("/api")
//	srv.AddAuthRoute("GET", "/json", "user", <callback>)
//	srv.AddAuthRoute("GET", "/json", []string{"user", "api"}, <callback>)
func (srv *Server) AddAuthRoute(method string, pattern string, role interface{}, handler func(r *http.Request) *Response) *mux.Route {
	return srv.AddAuthRouteWithWriter(method, pattern, role, func(w http.ResponseWriter, r *http.Request) *Response {
		return handler(r)
	})
}

// AddRouteWithWriter adds a route with writer
func (srv *Server) AddRouteWithWriter(method string, pattern string, handler func(w http.ResponseWriter, r *http.Request) *Response) *mux.Route {
	handle := func(w http.ResponseWriter, r *http.Request) {
		response := handler(w, r)
		if response == nil {
			return
		}

		if response.ContentType != "" {
			w.Header().Set("Content-Type", response.ContentType)
		}

		if response.Status != http.StatusOK {
			// call writeHeader explicit only in case of an error
			w.WriteHeader(response.Status)
		}
		if response.Body != nil {
			_, err := w.Write(response.Body)
			if err != nil {
				log.Warning("Failed to write response body: ", err.Error())
			}
		}
	}
	return srv.router.HandleFunc(srv.baseURL+pattern, handle).Methods(method)
}

// AddRoute adds a route
//
//	srv := rest.MakeServer("/api")
//	srv.AddRoute("GET", "/json", <callback>)
func (srv *Server) AddRoute(method string, pattern string, handler func(r *http.Request) *Response) *mux.Route {
	return srv.AddRouteWithWriter(method, pattern, func(w http.ResponseWriter, r *http.Request) *Response {
		return handler(r)
	})
}

// NewServer create a new REST API handler
func NewServer(baseURL string, listen Listener, routes []Route) (*Server, error) {
	srv := MakeServer(baseURL)

	for _, route := range routes {
		if route.Role == "noauth" {
			srv.AddRoute(route.Method, route.Pattern, route.Handler)
		} else {
			srv.AddAuthRoute(route.Method, route.Pattern, route.Role, route.Handler)
		}
	}

	log.Info("Listening on: " + listen.Proto + ":" + listen.Address + " with URI: " + baseURL)
	listener, err := Listen(listen.Proto, listen.Address, listen.Group)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to create rest listener: %v", err)
	}

	srv.listener = listener
	return srv, nil
}

// Serve starts the REST API handler, please consider using the method AsyncServe instead
func (srv *Server) Serve() error {
	return http.Serve(srv.listener, srv.GetRouter())
}

// AsyncServe is an asynchronous method for Serve, the returned channel may be used in the select block at the end of an app to encounter problems if the REST server does not serve REST requests or if it could not be initialized
func (srv *Server) AsyncServe() chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Serve()
	}()
	return errChan
}

// AddHandlerWS adds handler to server which handles websocket-setup
func (srv *Server) AddHandlerWS(method string, pattern string, wsRole string, handler func(r *http.Request, ws *websocket.Conn) uint16) *mux.Route {
	return srv.AddAuthSocket(method, pattern, wsRole, handler)
}

// AddAuthSocket upgrades and authorizes connection
func (srv *Server) AddAuthSocket(method string, pattern string, wsRole interface{}, handler func(r *http.Request, ws *websocket.Conn) uint16) *mux.Route {
	handle := func(w http.ResponseWriter, r *http.Request) {
		var upgrader websocket.Upgrader
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warningf("failed to upgrade websocket.Upgrader: %v", err)
			return
		}
		defer conn.Close()

		err = checkAuth(wsRole, conn)
		if err != nil {
			log.Warningf("failed to check authentication message: %v", err)
			err2 := sendClose(conn, websocket.ClosePolicyViolation, wsTimeout)
			log.Warningf("failed send via websocket: %v", err2)
			return
		}

		code := handler(r, conn)
		err = sendClose(conn, code, wsTimeout)
		log.Warningf("failed send handler: %v", err)
	}
	return srv.router.HandleFunc(srv.baseURL+pattern, handle).Methods(method)
}

func checkAuth(role interface{}, conn *websocket.Conn) error {
	msgtype, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("unable to read message: %v", err)
	}
	if msgtype != websocket.TextMessage {
		return errors.New("invalid authentication message")
	}

	return CheckAuth(role, string(msg))
}

func sendClose(conn *websocket.Conn, code uint16, wsTimeout time.Duration) error {
	var body [2]byte
	binary.BigEndian.PutUint16(body[:], code)
	err := conn.WriteControl(websocket.CloseMessage, body[:], time.Now().Add(wsTimeout))
	if err != nil {
		return fmt.Errorf("failed to write close to WebSocket: %v ", err)
	}
	return nil
}
