/*
 * device package - device.go
 * Copyright (c) 2018, TQ-Systems GmbH
 * All rights reserved. For further information see LICENSE.
 */

package device

import (
	"net"
	"os"

	"github.com/tq-systems/public-go-utils/config"
	"github.com/tq-systems/public-go-utils/log"
	"github.com/vishvananda/netlink"
)

//go:generate mockgen -destination=../mocks/device/mock_device.go -package=device github.com/tq-systems/public-go-utils/device Info

/*
#include <deviceinfo.h>
#cgo LDFLAGS: -ldeviceinfo
*/
import "C"

const (
	configFileTimezone  = "/cfglog/system/timezone.json"
	flagFileInvalidTime = "/run/em/system/time-invalid"
)

// GetDeviceSerial returns serial of the device
func GetDeviceSerial() string {
	return C.GoString(C.deviceinfo_get_serial_str())
}

// GetFirmwareVersion returns firmware version of the device
func GetFirmwareVersion() string {
	return C.GoString(C.deviceinfo_get_firmware_version_str())
}

// GetHardwareRevision returns hardware revision of the device
func GetHardwareRevision() string {
	return C.GoString(C.deviceinfo_get_hardware_revision_str())
}

// GetDeviceMac returns mac address of the device
func GetDeviceMac() string {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Error("Failed to get MAC address of eth0")
		return ""
	}
	return iface.HardwareAddr.String()
}

// GetDeviceIP returns IP address of the device
func GetDeviceIP() string {
	var linkip *net.IPNet

	link, _ := netlink.LinkByName("eth0")
	addrlist, _ := netlink.AddrList(link, netlink.FAMILY_V4)
	for _, addr := range addrlist {
		if addr.Scope == int(netlink.SCOPE_UNIVERSE) {
			// Return DHCP/static IP
			return addr.IPNet.IP.String()
		}
		// Auto
		linkip = addr.IPNet
	}
	if linkip == nil {
		return ""
	}
	return linkip.IP.String()
}

// GetTimezone returns the configured timezone
func GetTimezone() (string, error) {
	var ret string
	err := config.ReadJSON(configFileTimezone, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// GetProductName returns the product name of the device
func GetProductName() string {
	return C.GoString(C.deviceinfo_get_product_name())
}

// Info is the interface for this package (will be device.Info from other packages).
// The other exported functions are preserved for backwards compatibility.
type Info interface {
	GetTimestampValidity() bool
	GetSerial() string
	GetFirmwareVersion() string
	GetHardwareRevision() string
	GetMac() string
	GetIP() string
	GetTimezone() (string, error)
	GetProductName() string
}

type deviceInfo struct{}

// NewInfo returns an interface for device information
func NewInfo() Info {
	return &deviceInfo{}
}

// GetTimestampValidity returns true if system time is valid, and false if it is invalid
func (d *deviceInfo) GetTimestampValidity() bool {
	// system time is valid if flag file does not exist
	_, err := os.Stat(flagFileInvalidTime)
	return os.IsNotExist(err)
}

// GetSerial returns the serial number of the device
func (d *deviceInfo) GetSerial() string {
	return GetDeviceSerial()
}

// GetFirmwareVersion returns firmware version of the device
func (d *deviceInfo) GetFirmwareVersion() string {
	return GetFirmwareVersion()
}

// GetHardwareRevision returns hardware revision of the device
func (d *deviceInfo) GetHardwareRevision() string {
	return GetHardwareRevision()
}

// GetMac returns mac address of the device
func (d *deviceInfo) GetMac() string {
	return GetDeviceMac()
}

// GetIP returns IP address of the device
func (d *deviceInfo) GetIP() string {
	return GetDeviceIP()
}

// GetTimezone returns the configured timezone
func (d *deviceInfo) GetTimezone() (string, error) {
	return GetTimezone()
}

// GetProductName returns the product name of the device
func (d *deviceInfo) GetProductName() string {
	return GetProductName()
}
