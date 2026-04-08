/*
 * device package - device.go
 * Copyright (c) 2018 - 2023 TQ-Systems GmbH <license@tq-group.com>, D-82229 Seefeld, Germany. All rights reserved.
 * Author: Christoph Krutz and the Energy Manager development team
 *
 * This software code contained herein is licensed under the terms and conditions of
 * the TQ-Systems Product Software License Agreement Version 1.0.1 or any later version.
 * You will find the corresponding license text in the LICENSE file.
 * In case of any license issues please contact license@tq-group.com.
 */
package device

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/tq-systems/public-go-utils/v3/config"
	"github.com/tq-systems/public-go-utils/v3/log"
	"github.com/vishvananda/netlink"
)

//go:generate mockgen --build_flags=--mod=mod -source=device.go -destination=../mocks/device/mock_device.go -package=device Info

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
	iface, err := net.InterfaceByName("br0")
	if err != nil {
		log.Error("Failed to get MAC address of br0")
		return ""
	}
	return iface.HardwareAddr.String()
}

// GetDeviceIP returns IP address of the device
func GetDeviceIP() string {
	var linkip *net.IPNet

	link, _ := netlink.LinkByName("br0")
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

// GetDeviceType returns the device type of the device (equivalent to hardware type)
func GetDeviceType() string {
	return C.GoString(C.deviceinfo_get_device_type_str())
}

// Info is the interface for this package (will be device.Info from other packages).
// The other exported functions are preserved for backwards compatibility.
type Info interface {
	GetTimestampValidity() bool
	GetSerial() string
	GetHwType() string                              // Get hwtype from uboot environment variable 'hwtype'. If not set or not returned by libdeviceinfo, it returns empty string.
	GetRaucCompatible() (RaucCompatibleInfo, error) // GetRaucCompatible returns the rauc compatible string of the device (/etc/rauc/system.conf), parsed into its components. The parsing is done only once and cached for subsequent calls.
	GetFirmwareVersion() string
	GetHardwareRevision() string
	GetMac() string
	GetIP() string
	GetTimezone() (string, error)
	GetProductName() string
	GetDeviceType() string
}

type deviceInfo struct {
	raucCompatibleCache *RaucCompatibleInfo
	raucCompatibleInit  sync.Once // ensures that the rauc compatible string is parsed only once and cached for subsequent calls
}

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

// Get hwtype from uboot environment variable 'hwtype'. If not set or not returned by libdeviceinfo, it returns empty string.
func (d *deviceInfo) GetHwType() string {
	return C.GoString(C.deviceinfo_get_hwtype_str())
}

// RaucCompatibleInfo holds the parsed information from the rauc compatible string, as well as the raw string and any parsing error that occurred.
type RaucCompatibleInfo struct {
	RawString        string
	ParsingError     error
	BundleMachine    string
	BundleCompatible string
	BundleVersion    string
	SpecVersion      uint
}

// parseRaucCompatible parses a raw RAUC compatible string into RaucCompatibleInfo.
// Parsing failures are reported only via RaucCompatibleInfo.ParsingError, so callers can always
// use the complete parsed state (including RawString and partial fields).
func parseRaucCompatible(raw string) RaucCompatibleInfo {
	// The rauc compatible string in /etc/rauc/system.conf is specified in project meta-emos
	// in recipes-devtools/emit/files/emit in the function get_compatible_config().
	// Format: <bundleMachine>/<bundleCompatible>/<specVersion>/<bundleVersion> for specVersion 1
	var info RaucCompatibleInfo
	info.RawString = raw

	// Check if the string is empty, eg. if libdeviceinfo returns empty string or null
	if raw == "" {
		info.ParsingError = fmt.Errorf("rauc compatible string is empty or null")
		return info
	}

	{
		// Check if at least version can be parsed:
		// ^[^/]+/[^/]+/(\d+)/
		// ^                   Start of string
		//  [^/]+              One or more characters that are not '/'
		//       /             A '/'
		//        [^/]+        One or more characters that are not '/'
		//             /       A '/'
		//              (\d+)  One or more digits (the spec version) - captured
		//                   / A '/', end of the part we want to check for now
		pattern := regexp.MustCompile(`^[^/]+/[^/]+/(\d+)/`)
		matches := pattern.FindStringSubmatch(raw)
		if matches == nil {
			info.ParsingError = fmt.Errorf("invalid format: does not match expected pattern")
			return info
		}

		// Extract and store spec version
		specVersion, err := strconv.Atoi(matches[1])
		if err != nil {
			info.ParsingError = fmt.Errorf("failed to parse spec version: %w", err)
			return info
		}
		info.SpecVersion = uint(specVersion)
	}

	// Parse fields based on spec version
	switch info.SpecVersion {
	case 1:
		// V1 format: bundleMachine/bundleCompatible/1/bundleVersion
		v1Pattern := regexp.MustCompile(`^([^/]+)/([^/]+)/1/([^/]+)$`)
		v1Matches := v1Pattern.FindStringSubmatch(raw)
		if v1Matches == nil {
			info.ParsingError = fmt.Errorf("invalid V1 format")
			return info
		}
		info.BundleMachine = v1Matches[1]
		info.BundleCompatible = v1Matches[2]
		info.BundleVersion = v1Matches[3]
	default:
		info.ParsingError = fmt.Errorf("unsupported spec version: %d", info.SpecVersion)
		return info
	}

	return info
}

// GetRaucCompatible returns the rauc compatible string of the device (/etc/rauc/system.conf),
// parsed into its components. The parsing is done only once and cached for subsequent calls.
func (d *deviceInfo) GetRaucCompatible() (RaucCompatibleInfo, error) {
	d.raucCompatibleInit.Do(func() {
		raw := C.GoString(C.deviceinfo_get_rauc_compatible_str())
		parsed := parseRaucCompatible(raw)
		d.raucCompatibleCache = &parsed
	})

	return *d.raucCompatibleCache, d.raucCompatibleCache.ParsingError
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

// GetDeviceType returns the device type of the device (equivalent to hardware type)
func (d *deviceInfo) GetDeviceType() string {
	return GetDeviceType()
}
