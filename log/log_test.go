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
	"log/syslog"
	"strings"
	"testing"
	"time"

	"github.com/tq-systems/public-go-utils/outputcapturer"
)

func TestLoglevel(t *testing.T) {

	testLoglevel("debug", Debug, t)
	testLoglevel("debug", Info, t)
	testLoglevel("debug", Warning, t)
	testLoglevel("debug", Error, t)
	testLoglevel("debug", Critical, t)

	testLoglevel("info", Info, t)
	testLoglevel("info", Warning, t)
	testLoglevel("info", Error, t)
	testLoglevel("info", Critical, t)

	testLoglevel("warning", Warning, t)
	testLoglevel("warning", Error, t)
	testLoglevel("warning", Critical, t)

	testLoglevel("error", Error, t)
	testLoglevel("error", Critical, t)

	testLoglevel("critical", Critical, t)
}

func TestNotLogged(t *testing.T) {
	testNotLogged("info", Debug, "Debug", t)

	testNotLogged("warning", Debug, "Debug", t)
	testNotLogged("warning", Info, "Info", t)

	testNotLogged("error", Debug, "Debug", t)
	testNotLogged("error", Info, "Info", t)
	testNotLogged("error", Warning, "Warning", t)

	testNotLogged("critical", Debug, "Debug", t)
	testNotLogged("critical", Info, "Info", t)
	testNotLogged("critical", Warning, "Warning", t)
	testNotLogged("critical", Error, "Error", t)
}

func TestLoglevelf(t *testing.T) {

	testLoglevelf("debug", Debugf, t)
	testLoglevelf("debug", Infof, t)
	testLoglevelf("debug", Warningf, t)
	testLoglevelf("debug", Errorf, t)
	testLoglevelf("debug", Criticalf, t)

	testLoglevelf("info", Infof, t)
	testLoglevelf("info", Warningf, t)
	testLoglevelf("info", Errorf, t)
	testLoglevelf("info", Criticalf, t)

	testLoglevelf("warning", Warningf, t)
	testLoglevelf("warning", Errorf, t)
	testLoglevelf("warning", Criticalf, t)

	testLoglevelf("error", Errorf, t)
	testLoglevelf("error", Criticalf, t)

	testLoglevelf("critical", Criticalf, t)
}

func TestNotLoggedf(t *testing.T) {
	testNotLoggedf("info", Debugf, "Debugf", t)

	testNotLoggedf("warning", Debugf, "Debugf", t)
	testNotLoggedf("warning", Infof, "Infof", t)

	testNotLoggedf("error", Debugf, "Debugf", t)
	testNotLoggedf("error", Infof, "Infof", t)
	testNotLoggedf("error", Warningf, "Warningf", t)

	testNotLoggedf("critical", Debugf, "Debugf", t)
	testNotLoggedf("critical", Infof, "Infof", t)
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

	outputcapturer.StartCaptureStderr(1)
	InitLogger(loglevelToLog, true)
	fn("Test")
	firstLine := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))[0]

	if !strings.Contains(firstLine, "Test") {
		t.Error("Expected contains: 'Test' but was: ", firstLine)
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

	outputcapturer.StartCaptureStderr(1)
	InitLogger(loglevelToLog, true)
	fn("Test")
	firstLine := outputcapturer.GetStderr(time.Duration(time.Millisecond * 500))[0]

	if !strings.Contains(firstLine, "Test") {
		t.Error("Expected contains: 'Test' but was: ", firstLine)
	}
}

func TestNotLoggedCuzOfConsoleLog(t *testing.T) {

	outputcapturer.StartCaptureStderr(1)
	InitLogger("debug", false)

	defer func() {
		if r := recover(); r != nil {
			// everything is fine
		} else {
			t.Error("Panic expected")
		}
	}()
	Info("Test")
	_ = outputcapturer.GetStderr(time.Duration(time.Second * 1))[0]
}

func TestLogLevelString(t *testing.T) {
	InitLogger("debug", false)
	if logWriter == nil {
		t.Error("logWriter should not be nil")
	}
	testLogLevelString(syslog.LOG_DEBUG, "debug", t)
	testLogLevelString(syslog.LOG_DEBUG, "Debug", t)

	testLogLevelString(syslog.LOG_INFO, "info", t)
	testLogLevelString(syslog.LOG_INFO, "Info", t)

	testLogLevelString(syslog.LOG_WARNING, "warning", t)
	testLogLevelString(syslog.LOG_WARNING, "Warning", t)

	testLogLevelString(syslog.LOG_ERR, "error", t)
	testLogLevelString(syslog.LOG_ERR, "Error", t)

	testLogLevelString(syslog.LOG_CRIT, "critical", t)
	testLogLevelString(syslog.LOG_CRIT, "Critical", t)
}

func testLogLevelString(exptectedLogLevel syslog.Priority, loglevelAsString string, t *testing.T) {
	InitLogger(loglevelAsString, false)
	if logLevel != exptectedLogLevel {
		t.Error("logLevel should: ", exptectedLogLevel, " but was:", logLevel)
	}

}
