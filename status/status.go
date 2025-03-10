/*
 * status package - status.go
 * Copyright (c) 2019 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Matthias Schiffer and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */

package status

//go:generate mockgen --build_flags=--mod=mod -destination=../mocks/status/mock_status.go -package=status github.com/tq-systems/public-go-utils/v3/status Handler

import (
	"encoding/json"
	"fmt"
	"os"

	dbus "github.com/godbus/dbus/v5"
)

// SystemStatus is an enum describing the global system state
type SystemStatus int

// SystemStatus enum definitons
// Status of a group (e.g. update) must be consecutive
const (
	StatusIdle SystemStatus = iota
	StatusRebooting
	StatusUpdateUploading
	StatusUpdateValidating
	StatusUpdateInstalling
	StatusUpdateFinalizing
	StatusBackupExport
	StatusBackupImport
)

const (
	// D-Bus names
	updaterAppServiceName      = "com.tq_group.tq_em.updater1"
	updaterAppPathName         = "/com/tq_group/tq_em/updater1"
	getSystemStatusServiceName = "com.tq_group.tq_em.updater1.SystemStatus.GetStatus"
	setSystemStatusServiceName = "com.tq_group.tq_em.updater1.SystemStatus.SetStatus"
)

// MarshalJSON is the custom marshalling implementation for SystemStatus
func (s SystemStatus) MarshalJSON() ([]byte, error) {
	statusStrings := map[SystemStatus]string{
		StatusIdle:             "idle",
		StatusRebooting:        "rebooting",
		StatusUpdateUploading:  "update-uploading",
		StatusUpdateValidating: "update-validating",
		StatusUpdateInstalling: "update-installing",
		StatusUpdateFinalizing: "update-finalizing",
		StatusBackupExport:     "backup-export",
		StatusBackupImport:     "backup-import",
	}

	return json.Marshal(statusStrings[s])
}

// Handler is the status handler interface
type Handler interface {
	IsBusy() (bool, error)
	GetStatus() (SystemStatus, error)
	GetSafeMode() bool
	SetStatus(SystemStatus) (bool, error)
	SetStatusIfIdle(newStatus SystemStatus) (bool, error)
}

type handler struct {
	dbusObject dbus.BusObject
}

// NewStatus returns an interface for status information
func NewStatus() (Handler, error) {
	handle := &handler{}

	// connect to system bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	// d-bus object for getting/setting system state
	handle.dbusObject = conn.Object(updaterAppServiceName, dbus.ObjectPath(updaterAppPathName))

	return handle, nil
}

// IsBusy returns true if the current system status is not idle
func (h *handler) IsBusy() (bool, error) {
	status, err := h.GetStatus()
	if err != nil {
		return status != StatusIdle, fmt.Errorf("unable to get status: %v", err)
	}
	return status != StatusIdle, nil
}

// GetSafeMode returns true if the device is in safe mode
func (h *handler) GetSafeMode() bool {
	return fileExists("/update/safe-mode")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// SetStatusIfIdle sets a new system status if current state is idle;
// if the system is currently busy, the status is unchanged and false is returned
func (h *handler) SetStatusIfIdle(newStatus SystemStatus) (bool, error) {

	if busy, err := h.IsBusy(); busy {

		return false, fmt.Errorf("unable to find out if busy: %v", err)
	}

	status, err := h.SetStatus(newStatus)
	if err != nil {
		return status, fmt.Errorf("unable to set status: %v", err)
	}
	return status, nil
}

// GetStatus returns the current system status
func (h *handler) GetStatus() (SystemStatus, error) {
	status := StatusIdle

	err := h.dbusObject.Call(getSystemStatusServiceName, 0).Store(&status)
	if err != nil {
		return status, fmt.Errorf("dbus error: %v", err)

	}

	return status, nil
}

// SetStatus tries to set a new system status;
// returns true if the new system status could be set, otherwise returns false
func (h *handler) SetStatus(newStatus SystemStatus) (bool, error) {
	success := false

	err := h.dbusObject.Call(setSystemStatusServiceName, 0, newStatus).Store(&success)
	if err != nil {
		return success, fmt.Errorf("dbus error: %v", err)
	}

	return success, nil
}
