/*
 * mqtt handler - mqtthandler_test.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

//nolint:misspell
package mqtt

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tq-systems/public-go-utils/v3/log"
	"github.com/tq-systems/public-go-utils/v3/mqtt/test"
)

const (
	MQTTBrokerHost = "127.0.0.1"
	MQTTBrokerPort = 1884
	topic          = "TOPIC"
)

var (
	MQTTBroker = fmt.Sprintf("tcp://%s:%d", MQTTBrokerHost, MQTTBrokerPort)
	mutex      = sync.Mutex{}
)

func TestMQTTPubSubWithReconnect(t *testing.T) {
	waitGroup := &sync.WaitGroup{}

	mqttMessages := make([]*test.Test, 0)

	log.InitLogger("debug", true)
	log.Info("---", t.Name(), "---")

	broker := startMQTTBroker(t)
	time.Sleep(1 * time.Second)

	clientPub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTPublisher")
	if err != nil {
		t.Fatal(err)
	}
	defer clientPub.Close()

	clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
	if err != nil {
		t.Fatal(err)
	}
	defer clientSub.Close()

	_, err = clientSub.Subscribe(topic, waitForMessage(t, waitGroup, &mqttMessages))
	if err != nil {
		t.Fatal(err)
	}

	testMessage := &test.Test{}

	// stop broker and disconnect from clients
	stopMQTTBroker(t, broker)

	// start broker again and have clients reconnect to it automatically
	broker = startMQTTBroker(t)
	defer stopMQTTBroker(t, broker)

	// Wait for reconnect
	time.Sleep(3 * time.Second)

	waitGroup.Add(1)
	err = clientPub.Publish(topic, 0, false, testMessage)
	assert.Nil(t, err)

	testMessage.MessageCounter++

	waitGroup.Add(1)
	err = clientPub.Publish(topic, 0, false, testMessage)
	assert.Nil(t, err)

	// The test succeeds if we do not timeout here
	err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
	assert.Nil(t, err)

	assert.Equal(t, len(mqttMessages), 2, "Not enough messages")
	assert.Equal(t, (*mqttMessages[0]).MessageCounter, uint64(0), "Wrong MessageCounter for first message")
	assert.Equal(t, (*mqttMessages[1]).MessageCounter, uint64(1), "Wrong MessageCounter for second message")
}

func TestBroker(t *testing.T) {
	broker := startMQTTBroker(t)
	defer stopMQTTBroker(t, broker)

	t.Run("Subscribe to broker", func(t *testing.T) {
		waitGroup := &sync.WaitGroup{}

		clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
		if err != nil {
			t.Fatal(err)
		}
		defer clientSub.Close()

		subscription, err := clientSub.Subscribe(topic, waitForMessage(t, waitGroup, nil))
		assert.Nil(t, err)
		subscription.Unsubscribe()
	})

	t.Run("Subscribe to broker with empty topic", func(t *testing.T) {
		waitGroup := &sync.WaitGroup{}

		clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
		if err != nil {
			t.Fatal(err)
		}
		defer clientSub.Close()

		waitGroup.Add(1)
		go func() {

			_, err := clientSub.Subscribe("", waitForMessage(t, waitGroup, nil))
			assert.NotNil(t, err)
			waitGroup.Done()
		}()

		// Wait for the subscription process to finish in a timely manner
		// The test succeeds if we do not timeout here
		err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
		assert.Nil(t, err)
	})

	t.Run("Subscribe to broker with nil callback", func(t *testing.T) {
		waitGroup := &sync.WaitGroup{}

		clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
		if err != nil {
			t.Error(err)
		}
		defer clientSub.Close()

		waitGroup.Add(1)
		go func() {

			_, err := clientSub.Subscribe(topic, nil)
			assert.NotNil(t, err)
			waitGroup.Done()
		}()

		// Wait for the subscription process to finish in a timely manner
		// The test succeeds if we do not timeout here
		err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
		assert.NoError(t, err)
	})

	t.Run("Publish to broker", func(t *testing.T) {
		waitGroup := &sync.WaitGroup{}

		clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
		if err != nil {
			t.Fatal(err)
		}
		defer clientSub.Close()

		subscription, err := clientSub.Subscribe(topic, waitForMessage(t, waitGroup, nil))
		assert.Nil(t, err)
		defer subscription.Unsubscribe()

		clientPub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTPublisher")
		if err != nil {
			t.Fatal(err)
		}
		defer clientPub.Close()

		waitGroup.Add(1)
		err = clientPub.PublishEmpty(topic, 0, false)
		assert.Nil(t, err)

		// The test succeeds if we do not timeout here
		err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
		assert.Nil(t, err)
	})

	t.Run("Unsubscribe from broker", func(t *testing.T) {
		waitGroup := &sync.WaitGroup{}

		mqttMessages := make([]*test.Test, 0)

		clientPub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTPublisher")
		if err != nil {
			t.Error(err)
		}
		defer clientPub.Close()

		clientSub, err := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
		if err != nil {
			t.Error(err)
		}
		defer clientSub.Close()

		// Subscribe to the same topic twice to get one message twice
		sub, err := clientSub.Subscribe(topic, waitForMessage(t, waitGroup, &mqttMessages))
		assert.Nil(t, err)
		_, err = clientSub.Subscribe(topic, waitForMessage(t, waitGroup, &mqttMessages))
		assert.Nil(t, err)

		testMessage := &test.Test{}

		waitGroup.Add(2)
		err = clientPub.Publish(topic, 0, false, testMessage)
		assert.Nil(t, err)

		// Wait for the message to be received twice
		err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
		assert.Nil(t, err)

		// After unsubscription we receive the message once
		sub.Unsubscribe()
		time.Sleep(1 * time.Second)

		testMessage.MessageCounter++

		waitGroup.Add(1)
		err = clientPub.Publish(topic, 0, false, testMessage)
		assert.Nil(t, err)

		// Wait for the message to be received once
		err = waitWithTimeout(waitGroup, time.Duration(time.Second*2))
		assert.Nil(t, err)

		// Assertions
		assert.Equal(t, len(mqttMessages), 3, "Wrong message count")
		assert.Equal(t, mqttMessages[0].MessageCounter, uint64(0), "Wrong MessageCounter for message 1")
		assert.Equal(t, mqttMessages[1].MessageCounter, uint64(0), "Wrong MessageCounter for message 2")
		assert.Equal(t, mqttMessages[2].MessageCounter, uint64(1), "Wrong MessageCounter for message 3")
	})
}

func waitForMessage(t *testing.T, waitGroup *sync.WaitGroup, mqttMessages *[]*test.Test) func(topic string, msg []byte) {
	callbackProtomessage := func(topic string, msg []byte) {
		defer waitGroup.Done()
		var tests = test.Test{}
		err := tests.UnmarshalVT(msg)
		if err != nil {
			t.Fatal(err)
		}
		if mqttMessages != nil {
			mutex.Lock()
			defer mutex.Unlock()
			*mqttMessages = append(*mqttMessages, &tests)
		}
	}
	return callbackProtomessage
}

func startMQTTBroker(t *testing.T) *exec.Cmd {
	cmd := exec.Command("mosquitto", "-p", fmt.Sprint(MQTTBrokerPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	return cmd
}

func stopMQTTBroker(t *testing.T, broker *exec.Cmd) {
	err := broker.Process.Kill()
	if err != nil {
		t.Error(err)
	}
	_, err = broker.Process.Wait()
	if err != nil {
		t.Error(err)
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns errror if timeout reached
func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return errors.New("Timeout") // timed out
	}
}
