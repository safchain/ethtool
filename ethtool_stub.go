//go:build !linux

/*
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

// The ethtool package aims to provide a library that provides easy access
// to the Linux SIOCETHTOOL ioctl operations. It can be used to retrieve information
// from a network device such as statistics, driver related information or even
// the peer of a VETH interface.
package ethtool

import (
	"fmt"
	"runtime"
	"time"
)

var errOSUnsupported = fmt.Errorf("ethtool: %s is not supported", runtime.GOOS)

const (
	ETHTOOL_GLINKSETTINGS = 0x4c // unix.ETHTOOL_GLINKSETTINGS
	ETHTOOL_SLINKSETTINGS = 0x4d // unix.ETHTOOL_SLINKSETTINGS
)

// DriverName returns the driver name of the given interface name.
func (e *Ethtool) DriverName(intf string) (string, error) {
	return "", errOSUnsupported
}

// BusInfo returns the bus information of the given interface name.
func (e *Ethtool) BusInfo(intf string) (string, error) {
	return "", errOSUnsupported
}

// ModuleEeprom returns Eeprom information of the given interface name.
func (e *Ethtool) ModuleEeprom(intf string) ([]byte, error) {
	return nil, errOSUnsupported
}

// ModuleEepromHex returns Eeprom information as hexadecimal string
func (e *Ethtool) ModuleEepromHex(intf string) (string, error) {
	return "", errOSUnsupported
}

// DriverInfo returns driver information of the given interface name.
func (e *Ethtool) DriverInfo(intf string) (DrvInfo, error) {
	return DrvInfo{}, errOSUnsupported
}

// GetIndir retrieves the indirection table of the given interface name.
func (e *Ethtool) GetIndir(intf string) (Indir, error) {
	return Indir{}, errOSUnsupported
}

// SetIndir sets the indirection table of the given interface from the SetIndir struct
func (e *Ethtool) SetIndir(intf string, setIndir SetIndir) (Indir, error) {
	return Indir{}, errOSUnsupported
}

// GetChannels returns the number of channels for the given interface name.
func (e *Ethtool) GetChannels(intf string) (Channels, error) {
	return Channels{}, errOSUnsupported
}

// SetChannels sets the number of channels for the given interface name and
// returns the new number of channels.
func (e *Ethtool) SetChannels(intf string, channels Channels) (Channels, error) {
	return Channels{}, errOSUnsupported
}

// GetCoalesce returns the coalesce config for the given interface name.
func (e *Ethtool) GetCoalesce(intf string) (Coalesce, error) {
	return Coalesce{}, errOSUnsupported
}

// SetCoalesce sets the coalesce config for the given interface name.
func (e *Ethtool) SetCoalesce(intf string, coalesce Coalesce) (Coalesce, error) {
	return Coalesce{}, errOSUnsupported
}

// GetTimestampingInformation returns the PTP timestamping information for the given interface name.
func (e *Ethtool) GetTimestampingInformation(intf string) (TimestampingInformation, error) {
	return TimestampingInformation{}, errOSUnsupported
}

// PermAddr returns permanent address of the given interface name.
func (e *Ethtool) PermAddr(intf string) (string, error) {
	return "", errOSUnsupported
}

// GetWakeOnLan returns the WoL config for the given interface name.
func (e *Ethtool) GetWakeOnLan(intf string) (WakeOnLan, error) {
	return WakeOnLan{}, errOSUnsupported
}

// SetWakeOnLan sets the WoL config for the given interface name and
// returns the new WoL config.
func (e *Ethtool) SetWakeOnLan(intf string, wol WakeOnLan) (WakeOnLan, error) {
	return WakeOnLan{}, errOSUnsupported
}

// GetRing retrieves ring parameters of the given interface name.
func (e *Ethtool) GetRing(intf string) (Ring, error) {
	return Ring{}, errOSUnsupported
}

// SetRing sets ring parameters of the given interface name.
func (e *Ethtool) SetRing(intf string, ring Ring) (Ring, error) {
	return Ring{}, errOSUnsupported
}

// GetPause retrieves pause parameters of the given interface name.
func (e *Ethtool) GetPause(intf string) (Pause, error) {
	return Pause{}, errOSUnsupported
}

// SetPause sets pause parameters of the given interface name.
func (e *Ethtool) SetPause(intf string, pause Pause) (Pause, error) {
	return Pause{}, errOSUnsupported
}

// FeatureState contains the state of a feature.
type FeatureState struct {
	Available    bool
	Requested    bool
	Active       bool
	NeverChanged bool
}

// FeatureNames shows supported features by their name.
func (e *Ethtool) FeatureNames(intf string) (map[string]uint, error) {
	return nil, errOSUnsupported
}

// Features retrieves features of the given interface name.
func (e *Ethtool) Features(intf string) (map[string]bool, error) {
	return nil, errOSUnsupported
}

// FeaturesWithState retrieves features of the given interface name,
// with extra flags to explain if they can be enabled
func (e *Ethtool) FeaturesWithState(intf string) (map[string]FeatureState, error) {
	return nil, errOSUnsupported
}

// Change requests a change in the given device's features.
func (e *Ethtool) Change(intf string, config map[string]bool) error {
	return errOSUnsupported
}

// PrivFlagsNames shows supported private flags by their name.
func (e *Ethtool) PrivFlagsNames(intf string) (map[string]uint, error) {
	return nil, errOSUnsupported
}

// PrivFlags retrieves private flags of the given interface name.
func (e *Ethtool) PrivFlags(intf string) (map[string]bool, error) {
	return nil, errOSUnsupported
}

// UpdatePrivFlags requests a change in the given device's private flags.
func (e *Ethtool) UpdatePrivFlags(intf string, config map[string]bool) error {
	return errOSUnsupported
}

// LinkState get the state of a link.
func (e *Ethtool) LinkState(intf string) (uint32, error) {
	return 0, errOSUnsupported
}

// Stats retrieves stats of the given interface name.
// This maintains backward compatibility with existing code.
func (e *Ethtool) Stats(intf string) (map[string]uint64, error) {
	return nil, errOSUnsupported
}

// StatsWithBuffer retrieves stats of the given interface name using pre-allocated buffers.
// This allows the caller to control where the large structures are allocated,
// which can be useful to avoid heap allocations in Go 1.24+.
func (e *Ethtool) StatsWithBuffer(intf string, gstringsPtr *EthtoolGStrings, statsPtr *EthtoolStats) (map[string]uint64, error) {
	return nil, errOSUnsupported
}

// Close closes the ethool handler
func (e *Ethtool) Close() {
}

// Identity the nic with blink duration, if not specify blink for 60 seconds
func (e *Ethtool) Identity(intf string, duration *time.Duration) error {
	return errOSUnsupported
}

// NewEthtool returns a new ethtool handler
func NewEthtool() (*Ethtool, error) {
	return &Ethtool{}, nil
}

// BusInfo returns bus information of the given interface name.
func BusInfo(intf string) (string, error) {
	e, err := NewEthtool()
	if err != nil {
		return "", err
	}
	defer e.Close()
	return e.BusInfo(intf)
}

// DriverName returns the driver name of the given interface name.
func DriverName(intf string) (string, error) {
	e, err := NewEthtool()
	if err != nil {
		return "", err
	}
	defer e.Close()
	return e.DriverName(intf)
}

// Stats retrieves stats of the given interface name.
func Stats(intf string) (map[string]uint64, error) {
	e, err := NewEthtool()
	if err != nil {
		return nil, err
	}
	defer e.Close()
	return e.Stats(intf)
}

// PermAddr returns permanent address of the given interface name.
func PermAddr(intf string) (string, error) {
	e, err := NewEthtool()
	if err != nil {
		return "", err
	}
	defer e.Close()
	return e.PermAddr(intf)
}

// Identity the nic with blink duration, if not specify blink infinity
func Identity(intf string, duration *time.Duration) error {
	e, err := NewEthtool()
	if err != nil {
		return err
	}
	defer e.Close()
	return e.Identity(intf, duration)
}

// SupportedLinkModes returns the names of the link modes supported by the interface.
func SupportedLinkModes(mask uint64) []string {
	return nil
}

// SupportedSpeed returns the maximum capacity of this interface.
func SupportedSpeed(mask uint64) uint64 {
	return 0
}
