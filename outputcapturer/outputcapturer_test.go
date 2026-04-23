/*
 * Copyright (c) 2021-2026 TQ-Systems GmbH <license@tq-group.com>, D-82229
 * Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 *
 * This software is licensed under the TQ-Systems Product Software License
 * Agreement Version 1.0.3 or any later version.
 * You can obtain a copy of the License Agreement in the TQS (TQ-Systems
 * Software Licenses) folder on the following website:
 * https://www.tq-group.com/en/support/downloads/tq-software-license-conditions/
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
