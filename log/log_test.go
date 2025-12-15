/*
 * logging package - log_test.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Marcel Matzat and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */
package log

import (
	"strings"
	"testing"
	"time"

	"github.com/tq-systems/public-go-utils/v3/outputcapturer"
)

func TestLoglevel(t *testing.T) {

	testLoglevel("debug", Debug, t)
	testLoglevel("debug", Info, t)
	testLoglevel("debug", Notice, t)
	testLoglevel("debug", Warning, t)
	testLoglevel("debug", Error, t)
	testLoglevel("debug", Critical, t)

	testLoglevel("info", Info, t)
	testLoglevel("info", Notice, t)
	testLoglevel("info", Warning, t)
	testLoglevel("info", Error, t)
	testLoglevel("info", Critical, t)

	testLoglevel("notice", Notice, t)
	testLoglevel("notice", Warning, t)
	testLoglevel("notice", Error, t)
	testLoglevel("notice", Critical, t)

	testLoglevel("warning", Warning, t)
	testLoglevel("warning", Error, t)
	testLoglevel("warning", Critical, t)

	testLoglevel("error", Error, t)
	testLoglevel("error", Critical, t)

	testLoglevel("critical", Critical, t)
}

func TestNotLogged(t *testing.T) {
	testNotLogged("info", Debug, "Debug", t)

	testNotLogged("notice", Debug, "Debug", t)
	testNotLogged("notice", Info, "Info", t)

	testNotLogged("warning", Debug, "Debug", t)
	testNotLogged("warning", Info, "Info", t)
	testNotLogged("warning", Notice, "Notice", t)

	testNotLogged("error", Debug, "Debug", t)
	testNotLogged("error", Info, "Info", t)
	testNotLogged("error", Notice, "Notice", t)
	testNotLogged("error", Warning, "Warning", t)

	testNotLogged("critical", Debug, "Debug", t)
	testNotLogged("critical", Info, "Info", t)
	testNotLogged("critical", Notice, "Notice", t)
	testNotLogged("critical", Warning, "Warning", t)
	testNotLogged("critical", Error, "Error", t)
}

func TestLoglevelf(t *testing.T) {

	testLoglevelf("debug", Debugf, t)
	testLoglevelf("debug", Infof, t)
	testLoglevelf("debug", Noticef, t)
	testLoglevelf("debug", Warningf, t)
	testLoglevelf("debug", Errorf, t)
	testLoglevelf("debug", Criticalf, t)

	testLoglevelf("info", Infof, t)
	testLoglevelf("info", Noticef, t)
	testLoglevelf("info", Warningf, t)
	testLoglevelf("info", Errorf, t)
	testLoglevelf("info", Criticalf, t)

	testLoglevelf("notice", Noticef, t)
	testLoglevelf("notice", Warningf, t)
	testLoglevelf("notice", Errorf, t)
	testLoglevelf("notice", Criticalf, t)

	testLoglevelf("warning", Warningf, t)
	testLoglevelf("warning", Errorf, t)
	testLoglevelf("warning", Criticalf, t)

	testLoglevelf("error", Errorf, t)
	testLoglevelf("error", Criticalf, t)

	testLoglevelf("critical", Criticalf, t)
}

func TestNotLoggedf(t *testing.T) {
	testNotLoggedf("info", Debugf, "Debugf", t)

	testNotLoggedf("notice", Debugf, "Debug", t)
	testNotLoggedf("notice", Infof, "Info", t)

	testNotLoggedf("warning", Debugf, "Debugf", t)
	testNotLoggedf("warning", Infof, "Infof", t)
	testNotLoggedf("warning", Noticef, "Noticef", t)

	testNotLoggedf("error", Debugf, "Debugf", t)
	testNotLoggedf("error", Infof, "Infof", t)
	testNotLoggedf("error", Noticef, "Noticef", t)
	testNotLoggedf("error", Warningf, "Warningf", t)

	testNotLoggedf("critical", Debugf, "Debugf", t)
	testNotLoggedf("critical", Infof, "Infof", t)
	testNotLoggedf("critical", Noticef, "Noticef", t)
	testNotLoggedf("critical", Warningf, "Warningf", t)
	testNotLoggedf("critical", Errorf, "Errorf", t)
}

func testNotLogged(loglevelToLog string, fn func(...interface{}), fnName string, t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			// everything is fine
		} else {
			t.Error("Panic expected for loglevel: ", loglevelToLog, " with function: ", fnName)
		}
	}()

	testLoglevel(loglevelToLog, fn, t)
}

func testLoglevel(loglevelToLog string, fn func(...interface{}), t *testing.T) {

	err := outputcapturer.StartCaptureStderr(1)
	if err != nil {
		t.Error(err)
	}
	InitLogger(loglevelToLog, true)
	fn("Test")
	output := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))
	if len(output) == 0 {
		t.Error("Output was empty.")
	} else {
		firstLine := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))[0]

		if !strings.Contains(firstLine, "Test") {
			t.Error("Expected contains: 'Test' but was: ", firstLine)
		}
	}
}

func testNotLoggedf(loglevelToLog string, fn func(string, ...interface{}), fnName string, t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			// everything is fine
		} else {
			t.Error("Panic expected for loglevel: ", loglevelToLog, " with function: ", fnName)
		}
	}()

	testLoglevelf(loglevelToLog, fn, t)
}

func testLoglevelf(loglevelToLog string, fn func(string, ...interface{}), t *testing.T) {

	err := outputcapturer.StartCaptureStderr(1)
	if err != nil {
		t.Error(err)
	}
	InitLogger(loglevelToLog, true)
	fn("Test")
	firstLine := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))[0]

	if !strings.Contains(firstLine, "Test") {
		t.Error("Expected contains: 'Test' but was: ", firstLine)
	}
}

func TestConfigurationChangeFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string, ...interface{})
		loglevel string
	}{
		{
			name:     "ConfigurationChangeSystem logs at notice level",
			fn:       ConfigurationChangeInternalControl,
			loglevel: "notice",
		},
		{
			name:     "ConfigurationChangeUser logs at notice level",
			fn:       ConfigurationChangeUser,
			loglevel: "notice",
		},
		{
			name:     "ConfigurationChangeDBUS logs at notice level",
			fn:       ConfigurationChangeExternalControl,
			loglevel: "notice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := outputcapturer.StartCaptureStderr(1)
			if err != nil {
				t.Error(err)
			}
			InitLogger(tt.loglevel, true)
			tt.fn("test message")

			output := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))
			if len(output) == 0 {
				t.Errorf("Expected log output but got none for loglevel '%s'", tt.loglevel)
				return
			}

			firstLine := output[0]
			if !strings.Contains(firstLine, "test message") {
				t.Errorf("Expected log to contain '%s' but got: '%s'", "test message", firstLine)
			}
		})
	}
}

func TestConfigurationChangeFunctionsNotLogged(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string, ...interface{})
		fnName   string
		loglevel string
	}{
		{
			name:     "ConfigurationChangeSystem does not log when loglevel is warning",
			fn:       ConfigurationChangeInternalControl,
			fnName:   "ConfigurationChangeInternalControl",
			loglevel: "warning",
		},
		{
			name:     "ConfigurationChangeUser does not log when loglevel is warning",
			fn:       ConfigurationChangeUser,
			fnName:   "ConfigurationChangeUser",
			loglevel: "warning",
		},
		{
			name:     "ConfigurationChangeDBUS does not log when loglevel is warning",
			fn:       ConfigurationChangeExternalControl,
			fnName:   "ConfigurationChangeExternalControl",
			loglevel: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					// Panic expected because no output is generated (timeout in GetStderr)
				} else {
					t.Errorf("Panic expected for loglevel: %s with function: %s (no log output expected)", tt.loglevel, tt.fnName)
				}
			}()

			err := outputcapturer.StartCaptureStderr(1)
			if err != nil {
				t.Error(err)
			}
			InitLogger(tt.loglevel, true)
			tt.fn("test message")

			// This should timeout and panic because nothing should be logged
			outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))
		})
	}
}
