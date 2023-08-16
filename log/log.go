/*
 * logging package - log.go
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
	"fmt"
	"log/syslog"
	"os"
)

var (
	logWriter *syslog.Writer
	logLevel  = syslog.LOG_WARNING
)

func parseLoglevel(level string) syslog.Priority {
	switch level {
	case "debug", "Debug":
		return syslog.LOG_DEBUG
	case "info", "Info":
		return syslog.LOG_INFO
	case "warning", "Warning":
		return syslog.LOG_WARNING
	case "error", "Error":
		return syslog.LOG_ERR
	default:
		return syslog.LOG_CRIT
	}
}

// InitLogger to setup loglevel and if the logger should log to console
func InitLogger(level string, logToConsole bool) {
	logLevel = parseLoglevel(level)
	if logToConsole {
		logWriter = nil
	} else {
		logWriter, _ = syslog.New(logLevel, "")
	}
}

// Debug prints a message with debug priority
func Debug(args ...interface{}) {
	logMessage(syslog.LOG_DEBUG, args...)
}

// Debugf formats a message with debug priority
func Debugf(format string, args ...interface{}) {
	logfMessage(syslog.LOG_DEBUG, format, args...)
}

// Info prints a message with info priority
func Info(args ...interface{}) {
	logMessage(syslog.LOG_INFO, args...)
}

// Infof formats a message with info priority
func Infof(format string, args ...interface{}) {
	logfMessage(syslog.LOG_INFO, format, args...)
}

// Warning prints a message with warning priority
func Warning(args ...interface{}) {
	logMessage(syslog.LOG_WARNING, args...)
}

// Warningf formats a message with warning priority
func Warningf(format string, args ...interface{}) {
	logfMessage(syslog.LOG_WARNING, format, args...)
}

// Error prints a message with error priority
func Error(args ...interface{}) {
	logMessage(syslog.LOG_ERR, args...)
}

// Errorf formats a message with error priority
func Errorf(format string, args ...interface{}) {
	logfMessage(syslog.LOG_ERR, format, args...)
}

// Critical prints a message with critial priority
func Critical(args ...interface{}) {
	logMessage(syslog.LOG_CRIT, args...)
}

// Criticalf formats a message with critial priority
func Criticalf(format string, args ...interface{}) {
	logfMessage(syslog.LOG_CRIT, format, args...)
}

// Fatal prints a message with critial priority and calls os.Exit(1)
func Fatal(args ...interface{}) {
	Critical(args...)
	os.Exit(1)
}

// Fatalf formats a message with critial priority and calls os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	Criticalf(format, args...)
	os.Exit(1)
}

// Panic prints a message with critial priority and calls panic()
func Panic(args ...interface{}) {
	s := fmt.Sprint(args...)
	Critical(s)
	panic(s)
}

// Panicf formats a message with critial priority and calls panic()
func Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	Critical(s)
	panic(s)
}

func logSyslog(priority syslog.Priority, message string) {
	if priority <= syslog.LOG_CRIT {
		logWriter.Crit(message)
	} else if priority <= syslog.LOG_ERR {
		logWriter.Err(message)
	} else if priority <= syslog.LOG_WARNING {
		logWriter.Warning(message)
	} else if priority <= syslog.LOG_INFO {
		logWriter.Info(message)
	} else {
		logWriter.Debug(message)
	}
}

func logMessage(priority syslog.Priority, args ...interface{}) {
	if priority > logLevel {
		return
	}
	if logWriter == nil {
		// Log to console
		fmt.Fprintln(os.Stderr, args...)
	} else {
		logSyslog(priority, fmt.Sprint(args...))
	}
}

func logfMessage(priority syslog.Priority, format string, args ...interface{}) {
	if priority > logLevel {
		return
	}
	if logWriter == nil {
		// Log to console
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	} else {
		logSyslog(priority, fmt.Sprintf(format, args...))
	}
}
