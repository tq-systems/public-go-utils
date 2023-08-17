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

	"github.com/gogo/protobuf/proto"

	"github.com/tq-systems/public-go-utils/log"
	"github.com/tq-systems/public-go-utils/mqtt/test"
)

const (
	MQTTBrokerHost = "127.0.0.1"
	MQTTBrokerPort = 1884
)

var (
	MQTTBroker   = fmt.Sprintf("tcp://%s:%d", MQTTBrokerHost, MQTTBrokerPort)
	waitGroup    sync.WaitGroup
	mqttMessages []test.Test
)

func TestMQTTPubSubWithReconnect(t *testing.T) {
	mqttMessages = make([]test.Test, 0)

	log.InitLogger("debug", true)
	log.Info("---", t.Name(), "---")

	broker := startMQTTBroker(t)
	time.Sleep(1 * time.Second)

	clientPub := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTPublisher")
	clientSub := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")

	_, err := clientSub.Subscribe("TOPIC", callbackProtomessage)
	if err != nil {
		t.Error(err)
	}

	testMessage := &test.Test{}

	err = broker.Process.Kill()
	if err != nil {
		t.Error(err)
	}
	_, err = broker.Process.Wait()
	if err != nil {
		t.Error(err)
	}

	broker = startMQTTBroker(t)
	defer func() {
		clientPub.Close()
		clientSub.Close()

		err = broker.Process.Kill()
		if err != nil {
			t.Error(err)
		}
		_, err = broker.Process.Wait()
		if err != nil {
			t.Error(err)
		}
	}()

	// Wait for reconnect
	time.Sleep(3 * time.Second)

	waitGroup.Add(1)
	if err := clientPub.Publish("TOPIC", 0, false, testMessage); err != nil {
		t.Error(err)
	}

	testMessage.MessageCounter++

	waitGroup.Add(1)
	if err := clientPub.Publish("TOPIC", 0, false, testMessage); err != nil {
		t.Error(err)
	}

	if err := waitWithTimeout(&waitGroup, time.Duration(time.Second*2)); err != nil {
		t.Error(err)
	}

	if len(mqttMessages) != 2 {
		t.Error("Not enough messages")
	}

	if mqttMessages[0].MessageCounter != 0 {
		t.Error("Wrong MessageCounter for first message")
	}

	if mqttMessages[1].MessageCounter != 1 {
		t.Error("Wrong MessageCounter for second message")
	}
}

func TestMQTTUnsubscribe(t *testing.T) {
	mqttMessages = make([]test.Test, 0)

	log.InitLogger("debug", true)
	log.Info("---", t.Name(), "---")

	broker := startMQTTBroker(t)
	defer func() {
		err := broker.Process.Kill()
		if err != nil {
			t.Error(err)
		}
		_, err = broker.Process.Wait()
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(1 * time.Second)

	clientPub := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTPublisher")
	defer clientPub.Close()
	clientSub := NewClient(MQTTBrokerHost, MQTTBrokerPort, "MQTTSubscriber")
	defer clientSub.Close()

	sub, _ := clientSub.Subscribe("TOPIC", callbackProtomessage)
	_, _ = clientSub.Subscribe("TOPIC", callbackProtomessage)

	testMessage := &test.Test{}

	waitGroup.Add(2)
	if err := clientPub.Publish("TOPIC", 0, false, testMessage); err != nil {
		t.Error(err)
	}

	if err := waitWithTimeout(&waitGroup, time.Duration(time.Second*2)); err != nil {
		t.Error(err)
	}

	sub.Unsubscribe()
	time.Sleep(1 * time.Second)

	testMessage.MessageCounter++

	waitGroup.Add(1)
	if err := clientPub.Publish("TOPIC", 0, false, testMessage); err != nil {
		t.Error(err)
	}

	if err := waitWithTimeout(&waitGroup, time.Duration(time.Second*2)); err != nil {
		t.Error(err)
	}

	if len(mqttMessages) != 3 {
		t.Error("Wrong message count")
	}

	if mqttMessages[0].MessageCounter != 0 {
		t.Error("Wrong MessageCounter for message 1")
	}
	if mqttMessages[1].MessageCounter != 0 {
		t.Error("Wrong MessageCounter for message 2")
	}
	if mqttMessages[2].MessageCounter != 1 {
		t.Error("Wrong MessageCounter for message 3")
	}
}

func callbackProtomessage(topic string, msg []byte) {
	defer waitGroup.Done()
	var tests = test.Test{}
	_ = proto.Unmarshal(msg, &tests)
	mqttMessages = append(mqttMessages, tests)
}

func startMQTTBroker(t *testing.T) *exec.Cmd {
	cmd := exec.Command("mosquitto", "-p", fmt.Sprint(MQTTBrokerPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	return cmd
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
