/*
 * outputcapturer package - outputcapturer_test.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

package outputcapturer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrorCapturer(t *testing.T) {
	err := StartCaptureStderr(1)
	assert.NoError(t, err)
	fmt.Fprintln(os.Stderr, "Test Error")
	firstLine := GetStderr(time.Duration(time.Second * 2))[0]
	if firstLine != "Test Error" {
		t.Error("Expected: 'Test Error' but was: ", firstLine)
	}
}

// Expecting panic
func TestErrorCapturerTimeout(t *testing.T) {
	err := StartCaptureStderr(2)
	assert.NoError(t, err)
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
