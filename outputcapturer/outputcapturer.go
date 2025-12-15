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
	mu       sync.Mutex
	stderrMu sync.Mutex
)

// StartCaptureStderr starts capturing count Stderr calls. If count = 0, all outputs are captured
func StartCaptureStderr(count int) error {
	mu.Lock()
	stdError = stdError[:0]
	wg = &sync.WaitGroup{}
	wg.Add(count)
	mu.Unlock()

	stderrMu.Lock()
	osStdErr := os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		stderrMu.Unlock()
		return fmt.Errorf("failed to create pipe: %v", err)
	}
	os.Stderr = writer
	stderrMu.Unlock()

	var counter = 0
	go func() {
		defer func() {
			stderrMu.Lock()
			os.Stderr = osStdErr
			stderrMu.Unlock()
			reader.Close()
			writer.Close()
		}()

		scanner := bufio.NewScanner(reader)
		if count > 0 {
			for scanner.Scan() {
				counter++
				line := scanner.Text()

				mu.Lock()
				stdError = append(stdError, line)
				if wg != nil {
					wg.Done()
				}
				mu.Unlock()

				// count = 0 -> endless
				if counter >= count {
					break
				}
			}
		}
	}()

	return nil
}

// GetStderr provides all captured inputs and blocks till count (StartCaptureStderr) output captured. Panics after timeout
func GetStderr(timeout time.Duration) []string {

	doneChannel := make(chan bool, 1)
	go func() {
		mu.Lock()
		currentWg := wg
		mu.Unlock()

		if currentWg != nil {
			currentWg.Wait()
		}
		doneChannel <- true
	}()

	select {
	case <-doneChannel:
		//everything is fine
	case <-time.NewTimer(timeout).C:
		panic("Timeout occurred")
	}

	mu.Lock()
	// Return a copy to avoid race conditions when caller reads the slice while another goroutine modifies stdError
	result := make([]string, len(stdError))
	copy(result, stdError)
	mu.Unlock()

	return result
}
