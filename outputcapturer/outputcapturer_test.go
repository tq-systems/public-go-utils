/*
 * Copyright (c) 2018, TQ-Systems GmbH
 * All rights reserved. For further information see LICENSE.txt
 * Marcel Matzat
 */

package outputcapturer

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestErrorCapturer(t *testing.T) {
	StartCaptureStderr(1)
	fmt.Fprintln(os.Stderr, "Test Error")
	firstLine := GetStderr(time.Duration(time.Second * 2))[0]
	if firstLine != "Test Error" {
		t.Error("Expected: 'Test Error' but was: ", firstLine)
	}
}

// Expecting panic
func TestErrorCapturerTimeout(t *testing.T) {
	StartCaptureStderr(2)
	fmt.Fprintln(os.Stderr, "Test Error")

	defer func() {
		if r := recover(); r != nil {
			// everything is fine
		} else {
			t.Error("Panic expected")
		}
	}()
	GetStderr(time.Duration(time.Second * 2))

}
