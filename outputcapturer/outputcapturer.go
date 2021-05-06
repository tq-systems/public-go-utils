/*
 * Copyright (c) 2018, TQ-Systems GmbH
 * All rights reserved. For further information see LICENSE.txt
 * Marcel Matzat
 */

// Package outputcapturer is only for testting!
package outputcapturer

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

var (
	stdError                 []string
	reader, writer, osStdErr *os.File
	wg                       *sync.WaitGroup
)

// StartCaptureStderr starts capturing count Stderr calls. If count = 0, all outputs are captured
func StartCaptureStderr(count int) {
	stdError = stdError[:0]
	osStdErr = os.Stderr
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Panic(err)
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
				wg.Done()
				line := scanner.Text()
				stdError = append(stdError, line)

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
		panic("Timeout occured")
	}

	return stdError
}
