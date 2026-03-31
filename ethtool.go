//go:build linux

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
	"bytes"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	ETHTOOL_GLINKSETTINGS = unix.ETHTOOL_GLINKSETTINGS // 0x4c
	ETHTOOL_SLINKSETTINGS = unix.ETHTOOL_SLINKSETTINGS // 0x4d
)

var (
	gstringsPool = sync.Pool{
		New: func() interface{} {
			// new() will allocate and zero-initialize the struct.
			// The large data array within ethtoolGStrings will be zeroed.
			return new(EthtoolGStrings)
		},
	}
	statsPool = sync.Pool{
		New: func() interface{} {
			// new() will allocate and zero-initialize the struct.
			// The large data array within ethtoolStats will be zeroed.
			return new(EthtoolStats)
		},
	}

	// Updated supportedCapabilities including modes from ethtool.h enum ethtool_link_mode_bit_indices
	supportedCapabilities = []struct {
		name  string
		mask  uint64 // Use uint64 to accommodate indices > 31
		speed uint64 // Speed in bps, 0 for non-speed modes
	}{
		// Existing entries (reordered slightly by bit index for clarity)
		{"10baseT_Half", unix.ETHTOOL_LINK_MODE_10baseT_Half_BIT, 10_000_000},        // 0
		{"10baseT_Full", unix.ETHTOOL_LINK_MODE_10baseT_Full_BIT, 10_000_000},        // 1
		{"100baseT_Half", unix.ETHTOOL_LINK_MODE_100baseT_Half_BIT, 100_000_000},     // 2
		{"100baseT_Full", unix.ETHTOOL_LINK_MODE_100baseT_Full_BIT, 100_000_000},     // 3
		{"1000baseT_Half", unix.ETHTOOL_LINK_MODE_1000baseT_Half_BIT, 1_000_000_000}, // 4
		{"1000baseT_Full", unix.ETHTOOL_LINK_MODE_1000baseT_Full_BIT, 1_000_000_000}, // 5
		// Newly added or re-confirmed based on full enum
		{"Autoneg", unix.ETHTOOL_LINK_MODE_Autoneg_BIT, 0},                                    // 6
		{"TP", unix.ETHTOOL_LINK_MODE_TP_BIT, 0},                                              // 7 (Twisted Pair port)
		{"AUI", unix.ETHTOOL_LINK_MODE_AUI_BIT, 0},                                            // 8 (AUI port)
		{"MII", unix.ETHTOOL_LINK_MODE_MII_BIT, 0},                                            // 9 (MII port)
		{"FIBRE", unix.ETHTOOL_LINK_MODE_FIBRE_BIT, 0},                                        // 10 (FIBRE port)
		{"BNC", unix.ETHTOOL_LINK_MODE_BNC_BIT, 0},                                            // 11 (BNC port)
		{"10000baseT_Full", unix.ETHTOOL_LINK_MODE_10000baseT_Full_BIT, 10_000_000_000},       // 12
		{"Pause", unix.ETHTOOL_LINK_MODE_Pause_BIT, 0},                                        // 13
		{"Asym_Pause", unix.ETHTOOL_LINK_MODE_Asym_Pause_BIT, 0},                              // 14
		{"2500baseX_Full", unix.ETHTOOL_LINK_MODE_2500baseX_Full_BIT, 2_500_000_000},          // 15
		{"Backplane", unix.ETHTOOL_LINK_MODE_Backplane_BIT, 0},                                // 16 (Backplane port)
		{"1000baseKX_Full", unix.ETHTOOL_LINK_MODE_1000baseKX_Full_BIT, 1_000_000_000},        // 17
		{"10000baseKX4_Full", unix.ETHTOOL_LINK_MODE_10000baseKX4_Full_BIT, 10_000_000_000},   // 18
		{"10000baseKR_Full", unix.ETHTOOL_LINK_MODE_10000baseKR_Full_BIT, 10_000_000_000},     // 19
		{"10000baseR_FEC", unix.ETHTOOL_LINK_MODE_10000baseR_FEC_BIT, 10_000_000_000},         // 20
		{"20000baseMLD2_Full", unix.ETHTOOL_LINK_MODE_20000baseMLD2_Full_BIT, 20_000_000_000}, // 21
		{"20000baseKR2_Full", unix.ETHTOOL_LINK_MODE_20000baseKR2_Full_BIT, 20_000_000_000},   // 22
		{"40000baseKR4_Full", unix.ETHTOOL_LINK_MODE_40000baseKR4_Full_BIT, 40_000_000_000},   // 23
		{"40000baseCR4_Full", unix.ETHTOOL_LINK_MODE_40000baseCR4_Full_BIT, 40_000_000_000},   // 24
		{"40000baseSR4_Full", unix.ETHTOOL_LINK_MODE_40000baseSR4_Full_BIT, 40_000_000_000},   // 25
		{"40000baseLR4_Full", unix.ETHTOOL_LINK_MODE_40000baseLR4_Full_BIT, 40_000_000_000},   // 26
		{"56000baseKR4_Full", unix.ETHTOOL_LINK_MODE_56000baseKR4_Full_BIT, 56_000_000_000},   // 27
		{"56000baseCR4_Full", unix.ETHTOOL_LINK_MODE_56000baseCR4_Full_BIT, 56_000_000_000},   // 28
		{"56000baseSR4_Full", unix.ETHTOOL_LINK_MODE_56000baseSR4_Full_BIT, 56_000_000_000},   // 29
		{"56000baseLR4_Full", unix.ETHTOOL_LINK_MODE_56000baseLR4_Full_BIT, 56_000_000_000},   // 30
		{"25000baseCR_Full", unix.ETHTOOL_LINK_MODE_25000baseCR_Full_BIT, 25_000_000_000},     // 31
		// Modes beyond bit 31 (require GLINKSETTINGS)
		{"25000baseKR_Full", unix.ETHTOOL_LINK_MODE_25000baseKR_Full_BIT, 25_000_000_000},                      // 32
		{"25000baseSR_Full", unix.ETHTOOL_LINK_MODE_25000baseSR_Full_BIT, 25_000_000_000},                      // 33
		{"50000baseCR2_Full", unix.ETHTOOL_LINK_MODE_50000baseCR2_Full_BIT, 50_000_000_000},                    // 34
		{"50000baseKR2_Full", unix.ETHTOOL_LINK_MODE_50000baseKR2_Full_BIT, 50_000_000_000},                    // 35
		{"100000baseKR4_Full", unix.ETHTOOL_LINK_MODE_100000baseKR4_Full_BIT, 100_000_000_000},                 // 36
		{"100000baseSR4_Full", unix.ETHTOOL_LINK_MODE_100000baseSR4_Full_BIT, 100_000_000_000},                 // 37
		{"100000baseCR4_Full", unix.ETHTOOL_LINK_MODE_100000baseCR4_Full_BIT, 100_000_000_000},                 // 38
		{"100000baseLR4_ER4_Full", unix.ETHTOOL_LINK_MODE_100000baseLR4_ER4_Full_BIT, 100_000_000_000},         // 39
		{"50000baseSR2_Full", unix.ETHTOOL_LINK_MODE_50000baseSR2_Full_BIT, 50_000_000_000},                    // 40
		{"1000baseX_Full", unix.ETHTOOL_LINK_MODE_1000baseX_Full_BIT, 1_000_000_000},                           // 41
		{"10000baseCR_Full", unix.ETHTOOL_LINK_MODE_10000baseCR_Full_BIT, 10_000_000_000},                      // 42
		{"10000baseSR_Full", unix.ETHTOOL_LINK_MODE_10000baseSR_Full_BIT, 10_000_000_000},                      // 43
		{"10000baseLR_Full", unix.ETHTOOL_LINK_MODE_10000baseLR_Full_BIT, 10_000_000_000},                      // 44
		{"10000baseLRM_Full", unix.ETHTOOL_LINK_MODE_10000baseLRM_Full_BIT, 10_000_000_000},                    // 45
		{"10000baseER_Full", unix.ETHTOOL_LINK_MODE_10000baseER_Full_BIT, 10_000_000_000},                      // 46
		{"2500baseT_Full", unix.ETHTOOL_LINK_MODE_2500baseT_Full_BIT, 2_500_000_000},                           // 47 (already present but reconfirmed)
		{"5000baseT_Full", unix.ETHTOOL_LINK_MODE_5000baseT_Full_BIT, 5_000_000_000},                           // 48
		{"FEC_NONE", unix.ETHTOOL_LINK_MODE_FEC_NONE_BIT, 0},                                                   // 49
		{"FEC_RS", unix.ETHTOOL_LINK_MODE_FEC_RS_BIT, 0},                                                       // 50 (Reed-Solomon FEC)
		{"FEC_BASER", unix.ETHTOOL_LINK_MODE_FEC_BASER_BIT, 0},                                                 // 51 (BaseR FEC)
		{"50000baseKR_Full", unix.ETHTOOL_LINK_MODE_50000baseKR_Full_BIT, 50_000_000_000},                      // 52
		{"50000baseSR_Full", unix.ETHTOOL_LINK_MODE_50000baseSR_Full_BIT, 50_000_000_000},                      // 53
		{"50000baseCR_Full", unix.ETHTOOL_LINK_MODE_50000baseCR_Full_BIT, 50_000_000_000},                      // 54
		{"50000baseLR_ER_FR_Full", unix.ETHTOOL_LINK_MODE_50000baseLR_ER_FR_Full_BIT, 50_000_000_000},          // 55
		{"50000baseDR_Full", unix.ETHTOOL_LINK_MODE_50000baseDR_Full_BIT, 50_000_000_000},                      // 56
		{"100000baseKR2_Full", unix.ETHTOOL_LINK_MODE_100000baseKR2_Full_BIT, 100_000_000_000},                 // 57
		{"100000baseSR2_Full", unix.ETHTOOL_LINK_MODE_100000baseSR2_Full_BIT, 100_000_000_000},                 // 58
		{"100000baseCR2_Full", unix.ETHTOOL_LINK_MODE_100000baseCR2_Full_BIT, 100_000_000_000},                 // 59
		{"100000baseLR2_ER2_FR2_Full", unix.ETHTOOL_LINK_MODE_100000baseLR2_ER2_FR2_Full_BIT, 100_000_000_000}, // 60
		{"100000baseDR2_Full", unix.ETHTOOL_LINK_MODE_100000baseDR2_Full_BIT, 100_000_000_000},                 // 61
		{"200000baseKR4_Full", unix.ETHTOOL_LINK_MODE_200000baseKR4_Full_BIT, 200_000_000_000},                 // 62
		{"200000baseSR4_Full", unix.ETHTOOL_LINK_MODE_200000baseSR4_Full_BIT, 200_000_000_000},                 // 63
		{"200000baseLR4_ER4_FR4_Full", unix.ETHTOOL_LINK_MODE_200000baseLR4_ER4_FR4_Full_BIT, 200_000_000_000}, // 64
		{"200000baseDR4_Full", unix.ETHTOOL_LINK_MODE_200000baseDR4_Full_BIT, 200_000_000_000},                 // 65
		{"200000baseCR4_Full", unix.ETHTOOL_LINK_MODE_200000baseCR4_Full_BIT, 200_000_000_000},                 // 66
		{"100baseT1_Full", unix.ETHTOOL_LINK_MODE_100baseT1_Full_BIT, 100_000_000},                             // 67 (Automotive/SPE)
		{"1000baseT1_Full", unix.ETHTOOL_LINK_MODE_1000baseT1_Full_BIT, 1_000_000_000},                         // 68 (Automotive/SPE)
		{"400000baseKR8_Full", unix.ETHTOOL_LINK_MODE_400000baseKR8_Full_BIT, 400_000_000_000},                 // 69
		{"400000baseSR8_Full", unix.ETHTOOL_LINK_MODE_400000baseSR8_Full_BIT, 400_000_000_000},                 // 70
		{"400000baseLR8_ER8_FR8_Full", unix.ETHTOOL_LINK_MODE_400000baseLR8_ER8_FR8_Full_BIT, 400_000_000_000}, // 71
		{"400000baseDR8_Full", unix.ETHTOOL_LINK_MODE_400000baseDR8_Full_BIT, 400_000_000_000},                 // 72
		{"400000baseCR8_Full", unix.ETHTOOL_LINK_MODE_400000baseCR8_Full_BIT, 400_000_000_000},                 // 73
		{"FEC_LLRS", unix.ETHTOOL_LINK_MODE_FEC_LLRS_BIT, 0},                                                   // 74 (Low Latency Reed-Solomon FEC)
		// PAM4 modes start here? Often indicated by lack of KR/CR/SR/LR or different naming
		{"100000baseKR_Full", unix.ETHTOOL_LINK_MODE_100000baseKR_Full_BIT, 100_000_000_000},                   // 75 (Likely 100GBASE-KR1)
		{"100000baseSR_Full", unix.ETHTOOL_LINK_MODE_100000baseSR_Full_BIT, 100_000_000_000},                   // 76 (Likely 100GBASE-SR1)
		{"100000baseLR_ER_FR_Full", unix.ETHTOOL_LINK_MODE_100000baseLR_ER_FR_Full_BIT, 100_000_000_000},       // 77 (Likely 100GBASE-LR1/ER1/FR1)
		{"100000baseCR_Full", unix.ETHTOOL_LINK_MODE_100000baseCR_Full_BIT, 100_000_000_000},                   // 78 (Likely 100GBASE-CR1)
		{"100000baseDR_Full", unix.ETHTOOL_LINK_MODE_100000baseDR_Full_BIT, 100_000_000_000},                   // 79
		{"200000baseKR2_Full", unix.ETHTOOL_LINK_MODE_200000baseKR2_Full_BIT, 200_000_000_000},                 // 80 (Likely 200GBASE-KR2)
		{"200000baseSR2_Full", unix.ETHTOOL_LINK_MODE_200000baseSR2_Full_BIT, 200_000_000_000},                 // 81 (Likely 200GBASE-SR2)
		{"200000baseLR2_ER2_FR2_Full", unix.ETHTOOL_LINK_MODE_200000baseLR2_ER2_FR2_Full_BIT, 200_000_000_000}, // 82 (Likely 200GBASE-LR2/etc)
		{"200000baseDR2_Full", unix.ETHTOOL_LINK_MODE_200000baseDR2_Full_BIT, 200_000_000_000},                 // 83
		{"200000baseCR2_Full", unix.ETHTOOL_LINK_MODE_200000baseCR2_Full_BIT, 200_000_000_000},                 // 84 (Likely 200GBASE-CR2)
		{"400000baseKR4_Full", unix.ETHTOOL_LINK_MODE_400000baseKR4_Full_BIT, 400_000_000_000},                 // 85 (Likely 400GBASE-KR4)
		{"400000baseSR4_Full", unix.ETHTOOL_LINK_MODE_400000baseSR4_Full_BIT, 400_000_000_000},                 // 86 (Likely 400GBASE-SR4)
		{"400000baseLR4_ER4_FR4_Full", unix.ETHTOOL_LINK_MODE_400000baseLR4_ER4_FR4_Full_BIT, 400_000_000_000}, // 87 (Likely 400GBASE-LR4/etc)
		{"400000baseDR4_Full", unix.ETHTOOL_LINK_MODE_400000baseDR4_Full_BIT, 400_000_000_000},                 // 88
		{"400000baseCR4_Full", unix.ETHTOOL_LINK_MODE_400000baseCR4_Full_BIT, 400_000_000_000},                 // 89 (Likely 400GBASE-CR4)
		{"100baseFX_Half", unix.ETHTOOL_LINK_MODE_100baseFX_Half_BIT, 100_000_000},                             // 90
		{"100baseFX_Full", unix.ETHTOOL_LINK_MODE_100baseFX_Full_BIT, 100_000_000},                             // 91
	}
)

type ifreq struct {
	ifr_name [IFNAMSIZ]byte
	ifr_data uintptr
}

// following structures comes from uapi/linux/ethtool.h
type ethtoolSsetInfo struct {
	cmd       uint32
	reserved  uint32
	sset_mask uint64
	data      [MAX_SSET_INFO]uint32
}

type ethtoolGetFeaturesBlock struct {
	available     uint32
	requested     uint32
	active        uint32
	never_changed uint32
}

type ethtoolGfeatures struct {
	cmd    uint32
	size   uint32
	blocks [MAX_FEATURE_BLOCKS]ethtoolGetFeaturesBlock
}

type ethtoolSetFeaturesBlock struct {
	valid     uint32
	requested uint32
}

type ethtoolSfeatures struct {
	cmd    uint32
	size   uint32
	blocks [MAX_FEATURE_BLOCKS]ethtoolSetFeaturesBlock
}

type ethtoolDrvInfo struct {
	cmd          uint32
	driver       [32]byte
	version      [32]byte
	fw_version   [32]byte
	bus_info     [32]byte
	erom_version [32]byte
	reserved2    [12]byte
	n_priv_flags uint32
	n_stats      uint32
	testinfo_len uint32
	eedump_len   uint32
	regdump_len  uint32
}

type ethtoolEeprom struct {
	cmd    uint32
	magic  uint32
	offset uint32
	len    uint32
	data   [EEPROM_LEN]byte
}

type ethtoolModInfo struct {
	cmd        uint32
	tpe        uint32
	eeprom_len uint32
	reserved   [8]uint32
}

type ethtoolLink struct {
	cmd  uint32
	data uint32
}

type ethtoolPermAddr struct {
	cmd  uint32
	size uint32
	data [PERMADDR_LEN]byte
}

// Convert zero-terminated array of chars (string in C) to a Go string.
func goString(s []byte) string {
	strEnd := bytes.IndexByte(s, 0)
	if strEnd == -1 {
		return string(s)
	}
	return string(s[:strEnd])
}

// DriverName returns the driver name of the given interface name.
func (e *Ethtool) DriverName(intf string) (string, error) {
	info, err := e.getDriverInfo(intf)
	if err != nil {
		return "", err
	}
	return goString(info.driver[:]), nil
}

// BusInfo returns the bus information of the given interface name.
func (e *Ethtool) BusInfo(intf string) (string, error) {
	info, err := e.getDriverInfo(intf)
	if err != nil {
		return "", err
	}
	return goString(info.bus_info[:]), nil
}

// ModuleEeprom returns Eeprom information of the given interface name.
func (e *Ethtool) ModuleEeprom(intf string) ([]byte, error) {
	eeprom, _, err := e.getModuleEeprom(intf)
	if err != nil {
		return nil, err
	}

	return eeprom.data[:eeprom.len], nil
}

// ModuleEepromHex returns Eeprom information as hexadecimal string
func (e *Ethtool) ModuleEepromHex(intf string) (string, error) {
	eeprom, _, err := e.getModuleEeprom(intf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(eeprom.data[:eeprom.len]), nil
}

// DriverInfo returns driver information of the given interface name.
func (e *Ethtool) DriverInfo(intf string) (DrvInfo, error) {
	i, err := e.getDriverInfo(intf)
	if err != nil {
		return DrvInfo{}, err
	}

	drvInfo := DrvInfo{
		Cmd:         i.cmd,
		Driver:      goString(i.driver[:]),
		Version:     goString(i.version[:]),
		FwVersion:   goString(i.fw_version[:]),
		BusInfo:     goString(i.bus_info[:]),
		EromVersion: goString(i.erom_version[:]),
		Reserved2:   goString(i.reserved2[:]),
		NPrivFlags:  i.n_priv_flags,
		NStats:      i.n_stats,
		TestInfoLen: i.testinfo_len,
		EedumpLen:   i.eedump_len,
		RegdumpLen:  i.regdump_len,
	}

	return drvInfo, nil
}

// GetIndir retrieves the indirection table of the given interface name.
func (e *Ethtool) GetIndir(intf string) (Indir, error) {
	indir, err := e.getIndir(intf)
	if err != nil {
		return Indir{}, err
	}

	return indir, nil
}

// SetIndir sets the indirection table of the given interface from the SetIndir struct
func (e *Ethtool) SetIndir(intf string, setIndir SetIndir) (Indir, error) {

	if setIndir.Equal != 0 && setIndir.Weight != nil {
		return Indir{}, fmt.Errorf("equal and weight options are mutually exclusive")
	}

	indir, err := e.GetIndir(intf)
	if err != nil {
		return Indir{}, err
	}

	newindir, err := e.setIndir(intf, indir, setIndir)
	if err != nil {
		return Indir{}, err
	}

	return newindir, nil
}

// GetChannels returns the number of channels for the given interface name.
func (e *Ethtool) GetChannels(intf string) (Channels, error) {
	channels, err := e.getChannels(intf)
	if err != nil {
		return Channels{}, err
	}

	return channels, nil
}

// SetChannels sets the number of channels for the given interface name and
// returns the new number of channels.
func (e *Ethtool) SetChannels(intf string, channels Channels) (Channels, error) {
	channels, err := e.setChannels(intf, channels)
	if err != nil {
		return Channels{}, err
	}

	return channels, nil
}

// GetCoalesce returns the coalesce config for the given interface name.
func (e *Ethtool) GetCoalesce(intf string) (Coalesce, error) {
	coalesce, err := e.getCoalesce(intf)
	if err != nil {
		return Coalesce{}, err
	}
	return coalesce, nil
}

// SetCoalesce sets the coalesce config for the given interface name.
func (e *Ethtool) SetCoalesce(intf string, coalesce Coalesce) (Coalesce, error) {
	coalesce, err := e.setCoalesce(intf, coalesce)
	if err != nil {
		return Coalesce{}, err
	}
	return coalesce, nil
}

// GetTimestampingInformation returns the PTP timestamping information for the given interface name.
func (e *Ethtool) GetTimestampingInformation(intf string) (TimestampingInformation, error) {
	ts, err := e.getTimestampingInformation(intf)
	if err != nil {
		return TimestampingInformation{}, err
	}
	return ts, nil
}

// PermAddr returns permanent address of the given interface name.
func (e *Ethtool) PermAddr(intf string) (string, error) {
	permAddr, err := e.getPermAddr(intf)
	if err != nil {
		return "", err
	}

	if permAddr.data[0] == 0 && permAddr.data[1] == 0 &&
		permAddr.data[2] == 0 && permAddr.data[3] == 0 &&
		permAddr.data[4] == 0 && permAddr.data[5] == 0 {
		return "", nil
	}

	return fmt.Sprintf("%x:%x:%x:%x:%x:%x",
		permAddr.data[0:1],
		permAddr.data[1:2],
		permAddr.data[2:3],
		permAddr.data[3:4],
		permAddr.data[4:5],
		permAddr.data[5:6],
	), nil
}

// GetWakeOnLan returns the WoL config for the given interface name.
func (e *Ethtool) GetWakeOnLan(intf string) (WakeOnLan, error) {
	wol := WakeOnLan{
		Cmd: ETHTOOL_GWOL,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&wol))); err != nil {
		return WakeOnLan{}, err
	}

	return wol, nil
}

// SetWakeOnLan sets the WoL config for the given interface name and
// returns the new WoL config.
func (e *Ethtool) SetWakeOnLan(intf string, wol WakeOnLan) (WakeOnLan, error) {
	wol.Cmd = ETHTOOL_SWOL

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&wol))); err != nil {
		return WakeOnLan{}, err
	}

	return wol, nil
}

func (e *Ethtool) ioctl(intf string, data uintptr) error {
	var name [IFNAMSIZ]byte
	copy(name[:], []byte(intf))

	ifr := ifreq{
		ifr_name: name,
		ifr_data: data,
	}

	_, _, ep := unix.Syscall(unix.SYS_IOCTL, uintptr(e.fd), SIOCETHTOOL, uintptr(unsafe.Pointer(&ifr)))
	if ep != 0 {
		return ep
	}

	return nil
}

func (e *Ethtool) getDriverInfo(intf string) (ethtoolDrvInfo, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return ethtoolDrvInfo{}, err
	}

	return drvinfo, nil
}

// parsing of do_grxfhindir from ethtool.c
func (e *Ethtool) getIndir(intf string) (Indir, error) {
	indir_head := Indir{
		Cmd:  ETHTOOL_GRXFHINDIR,
		Size: 0,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indir_head))); err != nil {
		return Indir{}, err
	}

	indir := Indir{
		Cmd:  ETHTOOL_GRXFHINDIR,
		Size: indir_head.Size,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indir))); err != nil {
		return Indir{}, err
	}

	return indir, nil
}

// parsing of do_srxfhindir from ethtool.c
func (e *Ethtool) setIndir(intf string, indir Indir, setIndir SetIndir) (Indir, error) {

	err := fillIndirTable(&indir.Size, indir.RingIndex[:], 0, 0, int(setIndir.Equal), setIndir.Weight, uint32(len(setIndir.Weight)))
	if err != nil {
		return Indir{}, err
	}

	if indir.Size == ETH_RXFH_INDIR_NO_CHANGE {
		indir.Size = MAX_INDIR_SIZE
		return indir, nil
	}

	indir.Cmd = ETHTOOL_SRXFHINDIR
	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indir))); err != nil {
		return Indir{}, err
	}

	return indir, nil
}

func fillIndirTable(indirSize *uint32, indir []uint32, rxfhindirDefault int,
	rxfhindirStart int, rxfhindirEqual int, rxfhindirWeight []uint32,
	numWeights uint32) error {

	switch {
	case rxfhindirEqual != 0:
		for i := uint32(0); i < *indirSize; i++ {
			indir[i] = uint32(rxfhindirStart) + (i % uint32(rxfhindirEqual))
		}
	case rxfhindirWeight != nil:
		var sum, partial uint32 = 0, 0
		var j, weight uint32
		for j = range numWeights {
			weight = rxfhindirWeight[j]
			sum += weight
		}

		if sum == 0 {
			return fmt.Errorf("at least one weight must be non-zero")
		}

		if sum > *indirSize {
			return fmt.Errorf("total weight exceeds the size of the indirection table")
		}

		j = ^uint32(0) // equivalent to -1 for unsigned
		for i := uint32(0); i < *indirSize; i++ {
			for i >= (*indirSize*partial)/sum {
				j++
				weight = rxfhindirWeight[j]
				partial += weight
			}
			indir[i] = uint32(rxfhindirStart) + j
		}
	case rxfhindirDefault != 0:
		*indirSize = 0
	default:
		*indirSize = ETH_RXFH_INDIR_NO_CHANGE
	}
	return nil
}

func (e *Ethtool) getChannels(intf string) (Channels, error) {
	channels := Channels{
		Cmd: ETHTOOL_GCHANNELS,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&channels))); err != nil {
		return Channels{}, err
	}

	return channels, nil
}

func (e *Ethtool) setChannels(intf string, channels Channels) (Channels, error) {
	channels.Cmd = ETHTOOL_SCHANNELS

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&channels))); err != nil {
		return Channels{}, err
	}

	return channels, nil
}

func (e *Ethtool) getCoalesce(intf string) (Coalesce, error) {
	coalesce := Coalesce{
		Cmd: ETHTOOL_GCOALESCE,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&coalesce))); err != nil {
		return Coalesce{}, err
	}

	return coalesce, nil
}

func (e *Ethtool) setCoalesce(intf string, coalesce Coalesce) (Coalesce, error) {
	coalesce.Cmd = ETHTOOL_SCOALESCE

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&coalesce))); err != nil {
		return Coalesce{}, err
	}

	return coalesce, nil
}

func (e *Ethtool) getTimestampingInformation(intf string) (TimestampingInformation, error) {
	ts := TimestampingInformation{
		Cmd: ETHTOOL_GET_TS_INFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ts))); err != nil {
		return TimestampingInformation{}, err
	}

	return ts, nil
}

func (e *Ethtool) getPermAddr(intf string) (ethtoolPermAddr, error) {
	permAddr := ethtoolPermAddr{
		cmd:  ETHTOOL_GPERMADDR,
		size: PERMADDR_LEN,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&permAddr))); err != nil {
		return ethtoolPermAddr{}, err
	}

	return permAddr, nil
}

func (e *Ethtool) getModuleEeprom(intf string) (ethtoolEeprom, ethtoolModInfo, error) {
	modInfo := ethtoolModInfo{
		cmd: ETHTOOL_GMODULEINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&modInfo))); err != nil {
		return ethtoolEeprom{}, ethtoolModInfo{}, err
	}

	eeprom := ethtoolEeprom{
		cmd:    ETHTOOL_GMODULEEEPROM,
		len:    modInfo.eeprom_len,
		offset: 0,
	}

	if modInfo.eeprom_len > EEPROM_LEN {
		return ethtoolEeprom{}, ethtoolModInfo{}, fmt.Errorf("eeprom size: %d is larger than buffer size: %d", modInfo.eeprom_len, EEPROM_LEN)
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&eeprom))); err != nil {
		return ethtoolEeprom{}, ethtoolModInfo{}, err
	}

	return eeprom, modInfo, nil
}

// GetRing retrieves ring parameters of the given interface name.
func (e *Ethtool) GetRing(intf string) (Ring, error) {
	ring := Ring{
		Cmd: ETHTOOL_GRINGPARAM,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ring))); err != nil {
		return Ring{}, err
	}

	return ring, nil
}

// SetRing sets ring parameters of the given interface name.
func (e *Ethtool) SetRing(intf string, ring Ring) (Ring, error) {
	ring.Cmd = ETHTOOL_SRINGPARAM

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ring))); err != nil {
		return Ring{}, err
	}

	return ring, nil
}

// GetPause retrieves pause parameters of the given interface name.
func (e *Ethtool) GetPause(intf string) (Pause, error) {
	pause := Pause{
		Cmd: ETHTOOL_GPAUSEPARAM,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&pause))); err != nil {
		return Pause{}, err
	}

	return pause, nil
}

// SetPause sets pause parameters of the given interface name.
func (e *Ethtool) SetPause(intf string, pause Pause) (Pause, error) {
	pause.Cmd = ETHTOOL_SPAUSEPARAM

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&pause))); err != nil {
		return Pause{}, err
	}

	return pause, nil
}

func isFeatureBitSet(blocks [MAX_FEATURE_BLOCKS]ethtoolGetFeaturesBlock, index uint) bool {
	return (blocks)[index/32].active&(1<<(index%32)) != 0
}

// FeatureState contains the state of a feature.
type FeatureState struct {
	Available    bool
	Requested    bool
	Active       bool
	NeverChanged bool
}

func getFeatureStateBits(blocks [MAX_FEATURE_BLOCKS]ethtoolGetFeaturesBlock, index uint) FeatureState {
	return FeatureState{
		Available:    (blocks)[index/32].available&(1<<(index%32)) != 0,
		Requested:    (blocks)[index/32].requested&(1<<(index%32)) != 0,
		Active:       (blocks)[index/32].active&(1<<(index%32)) != 0,
		NeverChanged: (blocks)[index/32].never_changed&(1<<(index%32)) != 0,
	}
}

func setFeatureBit(blocks *[MAX_FEATURE_BLOCKS]ethtoolSetFeaturesBlock, index uint, value bool) {
	blockIndex, bitIndex := index/32, index%32

	blocks[blockIndex].valid |= 1 << bitIndex

	if value {
		blocks[blockIndex].requested |= 1 << bitIndex
	} else {
		blocks[blockIndex].requested &= ^(1 << bitIndex)
	}
}

func (e *Ethtool) getNames(intf string, mask int) (map[string]uint, error) {
	ssetInfo := ethtoolSsetInfo{
		cmd:       ETHTOOL_GSSET_INFO,
		sset_mask: 1 << mask,
		data:      [MAX_SSET_INFO]uint32{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ssetInfo))); err != nil {
		return nil, err
	}

	/* we only read data on first index because single bit was set in sset_mask(0x10) */
	length := ssetInfo.data[0]
	if length == 0 {
		return map[string]uint{}, nil
	} else if length > MAX_GSTRINGS {
		return nil, fmt.Errorf("ethtool currently doesn't support more than %d entries, received %d", MAX_GSTRINGS, length)
	}

	gstrings := EthtoolGStrings{
		cmd:        ETHTOOL_GSTRINGS,
		string_set: uint32(mask),
		len:        length,
		data:       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&gstrings))); err != nil {
		return nil, err
	}

	result := make(map[string]uint)
	for i := 0; i != int(length); i++ {
		b := gstrings.data[i*ETH_GSTRING_LEN : i*ETH_GSTRING_LEN+ETH_GSTRING_LEN]
		key := goString(b)
		if key != "" {
			result[key] = uint(i)
		}
	}

	return result, nil
}

// FeatureNames shows supported features by their name.
func (e *Ethtool) FeatureNames(intf string) (map[string]uint, error) {
	return e.getNames(intf, ETH_SS_FEATURES)
}

// Features retrieves features of the given interface name.
func (e *Ethtool) Features(intf string) (map[string]bool, error) {
	names, err := e.FeatureNames(intf)
	if err != nil {
		return nil, err
	}

	length := uint32(len(names))
	if length == 0 {
		return map[string]bool{}, nil
	}

	features := ethtoolGfeatures{
		cmd:  ETHTOOL_GFEATURES,
		size: (length + 32 - 1) / 32,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&features))); err != nil {
		return nil, err
	}

	result := make(map[string]bool, length)
	for key, index := range names {
		result[key] = isFeatureBitSet(features.blocks, index)
	}

	return result, nil
}

// FeaturesWithState retrieves features of the given interface name,
// with extra flags to explain if they can be enabled
func (e *Ethtool) FeaturesWithState(intf string) (map[string]FeatureState, error) {
	names, err := e.FeatureNames(intf)
	if err != nil {
		return nil, err
	}

	length := uint32(len(names))
	if length == 0 {
		return map[string]FeatureState{}, nil
	}

	features := ethtoolGfeatures{
		cmd:  ETHTOOL_GFEATURES,
		size: (length + 32 - 1) / 32,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&features))); err != nil {
		return nil, err
	}

	var result = make(map[string]FeatureState, length)
	for key, index := range names {
		result[key] = getFeatureStateBits(features.blocks, index)
	}

	return result, nil
}

// Change requests a change in the given device's features.
func (e *Ethtool) Change(intf string, config map[string]bool) error {
	names, err := e.FeatureNames(intf)
	if err != nil {
		return err
	}

	length := uint32(len(names))

	features := ethtoolSfeatures{
		cmd:  ETHTOOL_SFEATURES,
		size: (length + 32 - 1) / 32,
	}

	for key, value := range config {
		if index, ok := names[key]; ok {
			setFeatureBit(&features.blocks, index, value)
		} else {
			return fmt.Errorf("unsupported feature %q", key)
		}
	}

	return e.ioctl(intf, uintptr(unsafe.Pointer(&features)))
}

// PrivFlagsNames shows supported private flags by their name.
func (e *Ethtool) PrivFlagsNames(intf string) (map[string]uint, error) {
	return e.getNames(intf, ETH_SS_PRIV_FLAGS)
}

// PrivFlags retrieves private flags of the given interface name.
func (e *Ethtool) PrivFlags(intf string) (map[string]bool, error) {
	names, err := e.PrivFlagsNames(intf)
	if err != nil {
		return nil, err
	}

	length := uint32(len(names))
	if length == 0 {
		return map[string]bool{}, nil
	}

	var val ethtoolLink
	val.cmd = ETHTOOL_GPFLAGS
	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&val))); err != nil {
		return nil, err
	}

	result := make(map[string]bool, length)
	for name, mask := range names {
		result[name] = val.data&(1<<mask) != 0
	}

	return result, nil
}

// UpdatePrivFlags requests a change in the given device's private flags.
func (e *Ethtool) UpdatePrivFlags(intf string, config map[string]bool) error {
	names, err := e.PrivFlagsNames(intf)
	if err != nil {
		return err
	}

	var curr ethtoolLink
	curr.cmd = ETHTOOL_GPFLAGS
	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&curr))); err != nil {
		return err
	}

	var update ethtoolLink
	update.cmd = ETHTOOL_SPFLAGS
	update.data = curr.data
	for name, value := range config {
		if index, ok := names[name]; ok {
			if value {
				update.data |= 1 << index
			} else {
				update.data &= ^(1 << index)
			}
		} else {
			return fmt.Errorf("unsupported priv flag %q", name)
		}
	}

	return e.ioctl(intf, uintptr(unsafe.Pointer(&update)))
}

// LinkState get the state of a link.
func (e *Ethtool) LinkState(intf string) (uint32, error) {
	x := ethtoolLink{
		cmd: ETHTOOL_GLINK,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&x))); err != nil {
		return 0, err
	}

	return x.data, nil
}

// Stats retrieves stats of the given interface name.
// This maintains backward compatibility with existing code.
func (e *Ethtool) Stats(intf string) (map[string]uint64, error) {
	// Create temporary buffers and delegate to StatsWithBuffer
	gstrings := gstringsPool.Get().(*EthtoolGStrings)
	stats := statsPool.Get().(*EthtoolStats)
	defer func() {
		gstringsPool.Put(gstrings)
		statsPool.Put(stats)
	}()

	return e.StatsWithBuffer(intf, gstrings, stats)
}

// StatsWithBuffer retrieves stats of the given interface name using pre-allocated buffers.
// This allows the caller to control where the large structures are allocated,
// which can be useful to avoid heap allocations in Go 1.24+.
func (e *Ethtool) StatsWithBuffer(intf string, gstringsPtr *EthtoolGStrings, statsPtr *EthtoolStats) (map[string]uint64, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return nil, err
	}

	if drvinfo.n_stats > MAX_GSTRINGS {
		return nil, fmt.Errorf("ethtool currently doesn't support more than %d entries, received %d", MAX_GSTRINGS, drvinfo.n_stats)
	}

	gstringsPtr.cmd = ETHTOOL_GSTRINGS
	gstringsPtr.string_set = ETH_SS_STATS
	gstringsPtr.len = drvinfo.n_stats

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(gstringsPtr))); err != nil {
		return nil, err
	}

	statsPtr.cmd = ETHTOOL_GSTATS
	statsPtr.n_stats = drvinfo.n_stats

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(statsPtr))); err != nil {
		return nil, err
	}

	result := make(map[string]uint64, drvinfo.n_stats)
	for i := 0; i != int(drvinfo.n_stats); i++ {
		b := gstringsPtr.data[i*ETH_GSTRING_LEN : (i+1)*ETH_GSTRING_LEN]

		strEnd := bytes.IndexByte(b, 0)
		if strEnd == -1 {
			strEnd = ETH_GSTRING_LEN
		}
		key := string(b[:strEnd])

		if len(key) != 0 {
			result[key] = statsPtr.data[i]
		}
	}

	return result, nil
}

// Close closes the ethool handler
func (e *Ethtool) Close() {
	unix.Close(e.fd)
}

// Identity the nic with blink duration, if not specify blink for 60 seconds
func (e *Ethtool) Identity(intf string, duration *time.Duration) error {
	dur := uint32(DEFAULT_BLINK_DURATION.Seconds())
	if duration != nil {
		dur = uint32(duration.Seconds())
	}
	return e.identity(intf, IdentityConf{Duration: dur})
}

func (e *Ethtool) identity(intf string, identity IdentityConf) error {
	identity.Cmd = ETHTOOL_PHYS_ID
	return e.ioctl(intf, uintptr(unsafe.Pointer(&identity)))
}

// NewEthtool returns a new ethtool handler
func NewEthtool() (*Ethtool, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM|unix.SOCK_CLOEXEC, unix.IPPROTO_IP)
	if err != nil {
		return nil, err
	}

	return &Ethtool{
		fd: int(fd),
	}, nil
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

func supportedSpeeds(mask uint64) (ret []struct {
	name  string
	mask  uint64
	speed uint64
}) {
	for _, mode := range supportedCapabilities {
		if mode.speed > 0 && ((1<<mode.mask)&mask) != 0 {
			ret = append(ret, mode)
		}
	}
	return ret
}

// SupportedLinkModes returns the names of the link modes supported by the interface.
func SupportedLinkModes(mask uint64) []string {
	var ret []string
	for _, mode := range supportedSpeeds(mask) {
		ret = append(ret, mode.name)
	}
	return ret
}

// SupportedSpeed returns the maximum capacity of this interface.
func SupportedSpeed(mask uint64) uint64 {
	var ret uint64
	for _, mode := range supportedSpeeds(mask) {
		if mode.speed > ret {
			ret = mode.speed
		}
	}
	return ret
}
