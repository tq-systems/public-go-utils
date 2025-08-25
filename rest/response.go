/*
 * REST utilities - response.go
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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tq-systems/public-go-utils/v3/log"
)

// Response structure
type Response struct {
	ContentType string
	Status      int
	Body        []byte
}

// NewEmptyResponseWithStatus returns empty response with status
func NewEmptyResponseWithStatus(status int) *Response {
	return &Response{
		"",
		status,
		nil,
	}
}

// NewEmptyResponse returns empty response
func NewEmptyResponse() *Response {
	return NewEmptyResponseWithStatus(http.StatusNoContent)
}

// InternalError returns internal error
func InternalError(err error) *Response {
	log.Error(err.Error())
	return NewEmptyResponseWithStatus(http.StatusInternalServerError)
}

// NewJSONResponseWithStatus returns JSON response with status
func NewJSONResponseWithStatus(status int, data interface{}) *Response {
	res, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("unable to marshal data: %v", err)
		return InternalError(err)
	}

	return &Response{
		"application/json",
		status,
		res,
	}
}

// NewJSONResponse returns JSON response
func NewJSONResponse(data interface{}) *Response {
	return NewJSONResponseWithStatus(http.StatusOK, data)
}

// NewErrorResponse creates an error response. It has a JSON body in the
// form that is required for our public APIs. code and details may be
// nil.
//
// msg must contain an error message that is a full sentence, explaining
// the error to the user. It must not include unnecessary implementation
// details.
func NewErrorResponse(status int, msg string, code *int, details any) *Response {
	return NewJSONResponseWithStatus(status, errorResp{
		Error: errorBody{
			Message: msg,
			Code:    code,
			Details: details,
		},
	})
}
