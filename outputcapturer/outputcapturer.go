/*
 * outputcapturer package - outputcapturer.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

// Package outputcapturer is only for testting!
package outputcapturer

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	stdError []string
	wg       *sync.WaitGroup
)

// StartCaptureStderr starts capturing count Stderr calls. If count = 0, all outputs are captured
func StartCaptureStderr(count int) error {
	stdError = stdError[:0]
	osStdErr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %v", err)
	}
	os.Stderr = writer
	wg = &sync.WaitGroup{}
	wg.Add(count)
	var counter = 0
	go func() {
		scanner := bufio.NewScanner(reader)
		if count > 0 {
			for scanner.Scan() {
				counter++
				line := scanner.Text()
				stdError = append(stdError, line)
				wg.Done()

				// count = 0 -> endless
				if counter >= count {
					break
				}
			}
		}
		os.Stderr = osStdErr
		reader.Close()
		writer.Close()
	}()

	return nil
}

// GetStderr provides all captured inputs and blocks till count (StartCaptureStderr) output captured. Panics after timeout
func GetStderr(timeout time.Duration) []string {

	doneChannel := make(chan bool, 1)
	go func() {
		wg.Wait()
		doneChannel <- true
	}()

	select {
	case <-doneChannel:
		//everything is fine
	case <-time.NewTimer(timeout).C:
		panic("Timeout occurred")
	}

	return stdError
}
