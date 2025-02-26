/*
 * mqtt handler - mqtthandler.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-vtproto_out=. --go-vtproto_opt=features=unmarshal test/test.proto
//go:generate mockgen --build_flags=--mod=mod -destination=../mocks/mqtt/mock_client.go -package=mqtt github.com/tq-systems/public-go-utils/v3/mqtt Client
//go:generate mockgen --build_flags=--mod=mod -destination=../mocks/mqtt/mock_subscription.go -package=mqtt github.com/tq-systems/public-go-utils/v3/mqtt Subscription

//nolint:misspell
package mqtt

/*
#cgo LDFLAGS: -lmosquitto
#include <mosquitto.h>
#include <stdlib.h>

static void on_connect_cb(struct mosquitto *mosq, void *userdata, int result) {
	void onConnect(struct mosquitto *mosq);
	onConnect(mosq);
}

static void on_disconnect_cb(struct mosquitto *mosq, void *userdata, int result) {
	void onDisconnect(struct mosquitto *mosq);
	onDisconnect(mosq);
}

static void on_publish_cb(struct mosquitto *mosq, void *userdata, int mid) {
	void onPubSub(struct mosquitto *mosq, int mid);
	onPubSub(mosq, mid);
}

static void on_subscribe_cb(struct mosquitto *mosq, void *userdata, int mid, int qos_count, const int *granted_qos) {
	void onPubSub(struct mosquitto *mosq, int mid);
	onPubSub(mosq, mid);
}

static void on_message_cb(struct mosquitto *mosq, void *userdata, const struct mosquitto_message *msg) {
	void onMessage(struct mosquitto *mosq, struct mosquitto_message *msg);
	onMessage(mosq, (struct mosquitto_message *)msg);
}

static void setup_callbacks(struct mosquitto *mosq) {
	mosquitto_connect_callback_set(mosq, on_connect_cb);
	mosquitto_disconnect_callback_set(mosq, on_disconnect_cb);
	mosquitto_publish_callback_set(mosq, on_publish_cb);
	mosquitto_subscribe_callback_set(mosq, on_subscribe_cb);
	mosquitto_message_callback_set(mosq, on_message_cb);
}

*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/tq-systems/public-go-utils/v3/log"

	"google.golang.org/protobuf/proto"
)

// A Callback is a function run for every message received for a subscribed topic
type Callback func(topic string, message []byte)

type subscription struct {
	client   *client
	topic    string
	callback Callback
}

type client struct {
	mosq             *C.struct_mosquitto
	subscriptions    map[*subscription]bool
	subscribedTopics map[string]int

	connected     bool
	connectedCond *sync.Cond

	// Synchronizes accesses to the subscriptions maps and the connected condition
	lock *sync.Mutex

	/* Synchronizes subscribe/unsubscribe or robust publish calls
	 *
	 * While we do not support calling Subscribe/Unsubscribe from
	 * different goroutines, there may be explicit subscriptions
	 * running at the same time as resubscriptions due to automatic
	 * reconnect. doSubLock ensures that only one such call is
	 * running at a time.
	 *
	 * Explicit subscriptions always wait for the subscription to
	 * finish. This is done by storing channels for such subscriptions
	 * in confirmWaiters and blocking until the channel is closed. This
	 * happens when the subscription is confirmed in onPubSub or
	 * brokerConfirmTimeout has passed.
	 *
	 * The same facilities are used for robust publish calls (publish with
	 * QoS >= 1) to wait for the publish to succeed.
	 */
	currentMsgLock *sync.Mutex
	confirmWaiters map[C.int]chan error
}

// A Subscription tracks a registered subscription and can be used to unsubscribe
type Subscription interface {
	Unsubscribe()
}

// A Client represents a connection to an MQTT broker
type Client interface {
	Subscribe(topic string, callback Callback) (Subscription, error)
	PublishRaw(topic string, qos byte, retain bool, message []byte) error
	PublishEmpty(topic string, qos byte, retain bool) error
	Publish(topic string, qos byte, retain bool, message proto.Message) error
	Close()
}

var (
	initialize sync.Once
	lock       sync.Mutex

	// ErrConfirmTimedOut indicates that the MQTT broker did not confirm an
	// action within brokerConfirmTimeout.
	ErrConfirmTimedOut   = fmt.Errorf("waiting for confirmation from the broker timed out")
	brokerConfirmTimeout = 5 * time.Second

	// Global map to hold references to clients, and allow lookup from C callbacks
	// Must only be accessed with lock held
	clients = make(map[*C.struct_mosquitto]*client)
)

func locked(m *sync.Mutex, f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

// NewClient opens a new connection to an MQTT broker
func NewClient(brokerAddress string, brokerPort int, clientID string) (Client, error) {
	initialize.Do(func() {
		C.mosquitto_lib_init()
	})

	client := &client{
		subscriptions:    make(map[*subscription]bool),
		subscribedTopics: make(map[string]int),
		lock:             &sync.Mutex{},
		connectedCond:    &sync.Cond{},
		currentMsgLock:   &sync.Mutex{},
		confirmWaiters:   make(map[C.int]chan error),
	}
	client.connectedCond.L = client.lock

	cClientID := C.CString(clientID)
	defer C.free(unsafe.Pointer(cClientID))
	client.mosq = C.mosquitto_new(cClientID, true, nil)

	C.setup_callbacks(client.mosq)

	cBrokerAddress := C.CString(brokerAddress)
	defer C.free(unsafe.Pointer(cBrokerAddress))
	if C.mosquitto_connect(client.mosq, cBrokerAddress, C.int(brokerPort), 10) != 0 {
		log.Error(fmt.Sprintf("Unable to connect to MQTT broker %s:%d", brokerAddress, brokerPort))
		C.mosquitto_destroy(client.mosq)
		return nil, fmt.Errorf("unable to connect to MQTT broker %s:%d", brokerAddress, brokerPort)
	}

	locked(&lock, func() {
		clients[client.mosq] = client
	})

	C.mosquitto_loop_start(client.mosq)

	locked(client.lock, func() {
		for !client.connected {
			client.connectedCond.Wait()
		}
	})

	return client, nil
}

func getClient(client *C.struct_mosquitto) *client {
	lock.Lock()
	defer lock.Unlock()
	return clients[client]
}

/* onConnect updates the "connected" field of a client and ensures that subscriptions are
 * restored after an automatic reconnect to the MQTT broker.
 */
//export onConnect
func onConnect(mosq *C.struct_mosquitto) {
	log.Debug("MQTT connection established")

	client := getClient(mosq)

	topics := make([]string, 0)

	locked(client.lock, func() {
		for topic := range client.subscribedTopics {
			topics = append(topics, topic)
		}
	})

	for _, topic := range topics {
		err := client.doSubscribe(topic, false)
		if err != nil {
			log.Errorf("failed to subscribe to topic %s: %v", topic, err)
		}
	}

	locked(client.lock, func() {
		client.connected = true
		client.connectedCond.Broadcast()
	})
}

/* onDisconnect updates the "connected" field of a client. A warning
 * message is printed if the disconnect is unexpetected (not caused by
 * our own Close() call)
 */
//export onDisconnect
func onDisconnect(mosq *C.struct_mosquitto) {
	client := getClient(mosq)

	locked(client.lock, func() {
		if client.mosq != nil {
			log.Warning("MQTT connection lost")
		} else {
			log.Debug("MQTT connection closed")
		}

		client.connected = false
		client.connectedCond.Broadcast()
	})
}

// onPubSub wakes up a waiting doSubscribe/PublishRaw when a subscription/robust publish is finished.
//
//export onPubSub
func onPubSub(mosq *C.struct_mosquitto, mid C.int) {
	client := getClient(mosq)

	locked(client.currentMsgLock, func() {
		if ch, ok := client.confirmWaiters[mid]; ok {
			close(ch)
			delete(client.confirmWaiters, mid)
		}
	})
}

// getCallback returns the list of callbacks matching a given topic
func (client *client) getCallbacks(messageTopic *C.char) []Callback {
	callbacks := make([]Callback, 0)

	locked(client.lock, func() {
		for sub := range client.subscriptions {
			subTopic := C.CString(sub.topic)
			defer C.free(unsafe.Pointer(subTopic))

			var matches C.bool
			C.mosquitto_topic_matches_sub(subTopic, messageTopic, &matches)
			if matches {
				callbacks = append(callbacks, sub.callback)
			}
		}
	})

	return callbacks
}

/* onMessage handles incoming messages and runs the corresponding callbacks.
 * The callbacks are run synchronously in the mosquitto thread, so they must
 * not block; all more complex processing should be run in Goroutines.
 * Still, the callbacks run concurrently with the Go main thread, so accesses
 * to common data structures always need to be synchronized.
 */
//export onMessage
func onMessage(mosq *C.struct_mosquitto, message *C.struct_mosquitto_message) {
	client := getClient(mosq)
	callbacks := client.getCallbacks(message.topic)

	topic := C.GoString(message.topic)
	payload := C.GoBytes(message.payload, message.payloadlen)

	for _, cb := range callbacks {
		cb(topic, payload)
	}
}

func (client *client) Close() {
	mosq := client.mosq

	if mosq == nil {
		return
	}

	locked(client.lock, func() {
		client.mosq = nil

		C.mosquitto_disconnect(mosq)

		for client.connected {
			client.connectedCond.Wait()
		}
	})

	C.mosquitto_loop_stop(mosq, C.bool(true))
	C.mosquitto_destroy(mosq)

	locked(&lock, func() {
		delete(clients, mosq)
	})
}

/* doSubscribe is the low-level subscription function. It directly calls
 * mosquitto_subscribe (while holding doSubLock), optionally waiting for
 * the subscription request to finish.
 *
 * If wait is true, but the subscription was not confirmed within
 * brokerConfirmTimeout, ErrConfirmTimedOut will be returned.
 *
 * Note: doSubscribe is the only function of this package that may run either
 * from Go (when called through Subscribe) or from the libmosquitto handler
 * thread (when called from onConnect to restore subscriptions).
 */
func (client *client) doSubscribe(topic string, wait bool) error {
	cTopic := C.CString(topic)
	defer C.free(unsafe.Pointer(cTopic))

	var err error
	var publishDone chan error
	locked(client.currentMsgLock, func() {
		var currentSub C.int
		ret := C.mosquitto_subscribe(client.mosq, &currentSub, cTopic, 2)
		if ret != 0 {
			err = errors.New("Subscription of topic '" + topic + "' failed")
			return
		}
		if wait {
			// Only automatic resubscriptions set wait to false, and we do not want
			// to disturb explicit subscriptions running at the same time.
			publishDone = client.initConfirmWaiter(currentSub)
		}
	})
	if err == nil && publishDone != nil {
		err = <-publishDone
	}

	if err != nil {
		log.Debug("Subscribed on topic: " + topic)
	}

	return err
}

/* Subscribe adds a subscription for a topic (or topic pattern), running the given callback
 * for each message matching the topic. The locking found in the functions only synchoronizes
 * the thread/gorouting running Subscribe against the libmosquitto handler thread. The user of
 * the mqtt-handler must ensure that there are no concurrent calls to Subscribe/Unsubscribe.
 */
func (client *client) Subscribe(topic string, callback Callback) (Subscription, error) {
	if topic == "" || callback == nil {
		return nil, errors.New("error during Subscription: empty topic or nil callback not allowed")
	}
	sub := &subscription{
		client:   client,
		topic:    topic,
		callback: callback,
	}

	needSub := true

	locked(client.lock, func() {
		client.subscriptions[sub] = true
		client.subscribedTopics[topic]++

		if client.subscribedTopics[topic] > 1 {
			needSub = false
		}
	})

	if needSub {
		err := client.doSubscribe(topic, true)
		if err != nil {
			locked(client.lock, func() {
				delete(client.subscriptions, sub)
			})
			return nil, err
		}
	}

	return sub, nil
}

func (sub *subscription) Unsubscribe() {
	client := sub.client

	needUnsub := false

	locked(client.lock, func() {
		if client.subscriptions[sub] {
			delete(client.subscriptions, sub)
			client.subscribedTopics[sub.topic]--
			if client.subscribedTopics[sub.topic] == 0 {
				delete(client.subscribedTopics, sub.topic)
				needUnsub = true
			}
		}
	})

	if needUnsub {
		cTopic := C.CString(sub.topic)
		defer C.free(unsafe.Pointer(cTopic))

		locked(client.currentMsgLock, func() {
			C.mosquitto_unsubscribe(client.mosq, nil, cTopic)
		})
	}
}

// PublishRaw publishes a message to the MQTT broker.
//
// If qos is greater than 0, but the publication was not confirmed
// within brokerConfirmTimeout, ErrConfirmTimedOut will be returned.
func (client *client) PublishRaw(topic string, qos byte, retain bool, message []byte) error {
	cTopic := C.CString(topic)
	defer C.free(unsafe.Pointer(cTopic))

	var ptr unsafe.Pointer

	msglen := len(message)
	if msglen > 0 {
		ptr = unsafe.Pointer(&message[0])
	}

	var err error
	var publishDone chan error
	locked(client.currentMsgLock, func() {
		var currentMsg C.int
		ret := C.mosquitto_publish(client.mosq, &currentMsg, cTopic, C.int(msglen),
			ptr, C.int(qos), C.bool(retain))
		if ret != 0 {
			err = errors.New("failed to publish message")
			return
		}
		if qos > 0 {
			publishDone = client.initConfirmWaiter(currentMsg)
		}
	})
	if err == nil && publishDone != nil {
		err = <-publishDone
	}

	return err
}

func (client *client) PublishEmpty(topic string, qos byte, retain bool) error {
	return client.PublishRaw(topic, qos, retain, []byte{})
}

func (client *client) Publish(topic string, qos byte, retain bool, message proto.Message) error {
	marshalledProto, err := proto.Marshal(message)
	if err != nil {
		log.Panic(err.Error())
		return err
	}

	return client.PublishRaw(topic, qos, retain, marshalledProto)
}

// initConfirmWaiter adds a channel for mid to client.confirmWaiters
// and returns this channel. If the channel is still present in
// client.confirmWaiters after brokerConfirmTimeout, ErrConfirmTimedOut
// will be sent to it and it will be closed and removed again from
// client.confirmWaiters.
func (client *client) initConfirmWaiter(mid C.int) chan error {
	publishDone := make(chan error)
	client.confirmWaiters[mid] = publishDone
	time.AfterFunc(brokerConfirmTimeout, func() {
		client.currentMsgLock.Lock()
		defer client.currentMsgLock.Unlock()
		if ch, ok := client.confirmWaiters[mid]; ok {
			ch <- ErrConfirmTimedOut
			close(ch)
			delete(client.confirmWaiters, mid)
		}
	})
	return publishDone
}
