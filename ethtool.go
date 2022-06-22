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

// Package ethtool  aims to provide a library giving a simple access to the
// Linux SIOCETHTOOL ioctl operations. It can be used to retrieve informations
// from a network device like statistics, driver related informations or
// even the peer of a VETH interface.
package ethtool

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/bits"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// #include <stdlib.h>
import "C"

// Maximum size of an interface name
const (
	IFNAMSIZ = 16
)

// ioctl ethtool request
const (
	SIOCETHTOOL = 0x8946
)

type stringSet uint32

const (
	ETH_SS_TEST stringSet = iota
	ETH_SS_STATS
	ETH_SS_PRIV_FLAGS
	ETH_SS_NTUPLE_FILTERS
	ETH_SS_FEATURES
	ETH_SS_RSS_HASH_FUNCS
	ETH_SS_TUNABLES
	ETH_SS_PHY_STATS
	ETH_SS_PHY_TUNABLES
	ETH_SS_LINK_MODES
	ETH_SS_MSG_CLASSES
	ETH_SS_WOL_MODES
	ETH_SS_SOF_TIMESTAMPING
	ETH_SS_TS_TX_TYPES
	ETH_SS_TS_RX_FILTERS
	ETH_SS_UDP_TUNNEL_TYPES
	ETH_SS_STATS_STD
	ETH_SS_STATS_ETH_PHY
	ETH_SS_STATS_ETH_MAC
	ETH_SS_STATS_ETH_CTRL
	ETH_SS_STATS_RMON
)

// ethtool stats related constants.
const (
	ETH_GSTRING_LEN = 32

	// other CMDs from ethtool-copy.h of ethtool-3.5 package
	ETHTOOL_GSET     = 0x00000001 /* Get settings. */
	ETHTOOL_SSET     = 0x00000002 /* Set settings. */
	ETHTOOL_GDRVINFO = 0x00000003 /* Get driver info. */
	ETHTOOL_GREGS    = 0x00000004 /* Get NIC registers. */
	ETHTOOL_GWOL     = 0x00000005 /* Get wake-on-lan options. */
	ETHTOOL_SWOL     = 0x00000006 /* Set wake-on-lan options. */
	ETHTOOL_GMSGLVL  = 0x00000007 /* Get driver message level */
	ETHTOOL_SMSGLVL  = 0x00000008 /* Set driver msg level. */
	ETHTOOL_NWAY_RST = 0x00000009 /* Restart autonegotiation. */

	/* Get link status for host, i.e. whether the interface *and* the
	 * physical port (if there is one) are up (ethtool_value). */
	ETHTOOL_GLINK       = 0x0000000a
	ETHTOOL_GEEPROM     = 0x0000000b /* Get EEPROM data */
	ETHTOOL_SEEPROM     = 0x0000000c /* Set EEPROM data. */
	ETHTOOL_GCOALESCE   = 0x0000000e /* Get coalesce config */
	ETHTOOL_SCOALESCE   = 0x0000000f /* Set coalesce config. */
	ETHTOOL_GRINGPARAM  = 0x00000010 /* Get ring parameters */
	ETHTOOL_SRINGPARAM  = 0x00000011 /* Set ring parameters. */
	ETHTOOL_GPAUSEPARAM = 0x00000012 /* Get pause parameters */
	ETHTOOL_SPAUSEPARAM = 0x00000013 /* Set pause parameters. */
	ETHTOOL_GRXCSUM     = 0x00000014 /* Get RX hw csum enable (ethtool_value) */
	ETHTOOL_SRXCSUM     = 0x00000015 /* Set RX hw csum enable (ethtool_value) */
	ETHTOOL_GTXCSUM     = 0x00000016 /* Get TX hw csum enable (ethtool_value) */
	ETHTOOL_STXCSUM     = 0x00000017 /* Set TX hw csum enable (ethtool_value) */
	ETHTOOL_GSG         = 0x00000018 /* Get scatter-gather enable (ethtool_value) */
	ETHTOOL_SSG         = 0x00000019 /* Set scatter-gather enable (ethtool_value). */
	ETHTOOL_TEST        = 0x0000001a /* execute NIC self-test. */
	ETHTOOL_GSTRINGS    = 0x0000001b /* get specified string set */
	ETHTOOL_PHYS_ID     = 0x0000001c /* identify the NIC */
	ETHTOOL_GSTATS      = 0x0000001d /* get NIC-specific statistics */
	ETHTOOL_GTSO        = 0x0000001e /* Get TSO enable (ethtool_value) */
	ETHTOOL_STSO        = 0x0000001f /* Set TSO enable (ethtool_value) */
	ETHTOOL_GPERMADDR   = 0x00000020 /* Get permanent hardware address */
	ETHTOOL_GUFO        = 0x00000021 /* Get UFO enable (ethtool_value) */
	ETHTOOL_SUFO        = 0x00000022 /* Set UFO enable (ethtool_value) */
	ETHTOOL_GGSO        = 0x00000023 /* Get GSO enable (ethtool_value) */
	ETHTOOL_SGSO        = 0x00000024 /* Set GSO enable (ethtool_value) */
	ETHTOOL_GFLAGS      = 0x00000025 /* Get flags bitmap(ethtool_value) */
	ETHTOOL_SFLAGS      = 0x00000026 /* Set flags bitmap(ethtool_value) */
	ETHTOOL_GPFLAGS     = 0x00000027 /* Get driver-private flags bitmap */
	ETHTOOL_SPFLAGS     = 0x00000028 /* Set driver-private flags bitmap */

	ETHTOOL_GRXFH       = 0x00000029 /* Get RX flow hash configuration */
	ETHTOOL_SRXFH       = 0x0000002a /* Set RX flow hash configuration */
	ETHTOOL_GGRO        = 0x0000002b /* Get GRO enable (ethtool_value) */
	ETHTOOL_SGRO        = 0x0000002c /* Set GRO enable (ethtool_value) */
	ETHTOOL_GRXRINGS    = 0x0000002d /* Get RX rings available for LB */
	ETHTOOL_GRXCLSRLCNT = 0x0000002e /* Get RX class rule count */
	ETHTOOL_GRXCLSRULE  = 0x0000002f /* Get RX classification rule */
	ETHTOOL_GRXCLSRLALL = 0x00000030 /* Get all RX classification rule */
	ETHTOOL_SRXCLSRLDEL = 0x00000031 /* Delete RX classification rule */
	ETHTOOL_SRXCLSRLINS = 0x00000032 /* Insert RX classification rule */
	ETHTOOL_FLASHDEV    = 0x00000033 /* Flash firmware to device */
	ETHTOOL_RESET       = 0x00000034 /* Reset hardware */
	ETHTOOL_SRXNTUPLE   = 0x00000035 /* Add an n-tuple filter to device */
	ETHTOOL_GRXNTUPLE   = 0x00000036 /* deprecated */
	ETHTOOL_GSSET_INFO  = 0x00000037 /* Get string set info */
	ETHTOOL_GRXFHINDIR  = 0x00000038 /* Get RX flow hash indir'n table */
	ETHTOOL_SRXFHINDIR  = 0x00000039 /* Set RX flow hash indir'n table */

	ETHTOOL_GFEATURES     = 0x0000003a /* Get device offload settings */
	ETHTOOL_SFEATURES     = 0x0000003b /* Change device offload settings */
	ETHTOOL_GCHANNELS     = 0x0000003c /* Get no of channels */
	ETHTOOL_SCHANNELS     = 0x0000003d /* Set no of channels */
	ETHTOOL_SET_DUMP      = 0x0000003e /* Set dump settings */
	ETHTOOL_GET_DUMP_FLAG = 0x0000003f /* Get dump settings */
	ETHTOOL_GET_DUMP_DATA = 0x00000040 /* Get dump data */
	ETHTOOL_GET_TS_INFO   = 0x00000041 /* Get time stamping and PHC info */
	ETHTOOL_GMODULEINFO   = 0x00000042 /* Get plug-in module information */
	ETHTOOL_GMODULEEEPROM = 0x00000043 /* Get plug-in module eeprom */
	ETHTOOL_GEEE          = 0x00000044 /* Get EEE settings */
	ETHTOOL_SEEE          = 0x00000045 /* Set EEE settings */

	ETHTOOL_GRSSH     = 0x00000046 /* Get RX flow hash configuration */
	ETHTOOL_SRSSH     = 0x00000047 /* Set RX flow hash configuration */
	ETHTOOL_GTUNABLE  = 0x00000048 /* Get tunable configuration */
	ETHTOOL_STUNABLE  = 0x00000049 /* Set tunable configuration */
	ETHTOOL_GPHYSTATS = 0x0000004a /* get PHY-specific statistics */

	ETHTOOL_PERQUEUE = 0x0000004b /* Set per queue options */

	ETHTOOL_GLINKSETTINGS = 0x0000004c /* Get ethtool_link_settings */
	ETHTOOL_SLINKSETTINGS = 0x0000004d /* Set ethtool_link_settings */
	ETHTOOL_PHY_GTUNABLE  = 0x0000004e /* Get PHY tunable configuration */
	ETHTOOL_PHY_STUNABLE  = 0x0000004f /* Set PHY tunable configuration */
	ETHTOOL_GFECPARAM     = 0x00000050 /* Get FEC settings */
	ETHTOOL_SFECPARAM     = 0x00000051 /* Set FEC settings */
)

// MAX_GSTRINGS maximum number of stats entries that ethtool can
// retrieve currently.
const (
	MAX_GSTRINGS       = 32768
	MAX_FEATURE_BLOCKS = (MAX_GSTRINGS + 32 - 1) / 32
	EEPROM_LEN         = 640
	PERMADDR_LEN       = 32
	ETH_ALEN           = 6
)

type ifreq struct {
	ifr_name [IFNAMSIZ]byte
	ifr_data uintptr
}

// following structures comes from uapi/linux/ethtool.h
type ethtoolSsetInfo struct {
	cmd       uint32
	reserved  uint32
	sset_mask uint32
	data      uintptr
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

// DrvInfo contains driver information
// ethtool.h v3.5: struct ethtool_drvinfo
type DrvInfo struct {
	Cmd         uint32
	Driver      string
	Version     string
	FwVersion   string
	BusInfo     string
	EromVersion string
	Reserved2   string
	NPrivFlags  uint32
	NStats      uint32
	TestInfoLen uint32
	EedumpLen   uint32
	RegdumpLen  uint32
}

// Channels contains the number of channels for a given interface.
type Channels struct {
	Cmd           uint32
	MaxRx         uint32
	MaxTx         uint32
	MaxOther      uint32
	MaxCombined   uint32
	RxCount       uint32
	TxCount       uint32
	OtherCount    uint32
	CombinedCount uint32
}

// Coalesce is a coalesce config for an interface
type Coalesce struct {
	Cmd                      uint32
	RxCoalesceUsecs          uint32
	RxMaxCoalescedFrames     uint32
	RxCoalesceUsecsIrq       uint32
	RxMaxCoalescedFramesIrq  uint32
	TxCoalesceUsecs          uint32
	TxMaxCoalescedFrames     uint32
	TxCoalesceUsecsIrq       uint32
	TxMaxCoalescedFramesIrq  uint32
	StatsBlockCoalesceUsecs  uint32
	UseAdaptiveRxCoalesce    uint32
	UseAdaptiveTxCoalesce    uint32
	PktRateLow               uint32
	RxCoalesceUsecsLow       uint32
	RxMaxCoalescedFramesLow  uint32
	TxCoalesceUsecsLow       uint32
	TxMaxCoalescedFramesLow  uint32
	PktRateHigh              uint32
	RxCoalesceUsecsHigh      uint32
	RxMaxCoalescedFramesHigh uint32
	TxCoalesceUsecsHigh      uint32
	TxMaxCoalescedFramesHigh uint32
	RateSampleInterval       uint32
}

type ethtoolGStrings struct {
	cmd        uint32
	string_set uint32
	len        uint32
	data       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte
}

type ethtoolStats struct {
	cmd     uint32
	n_stats uint32
	data    [MAX_GSTRINGS]uint64
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

type (
	be16 = uint16
	be32 = uint32
)

type ethtoolRxfh struct {
	cmd         uint32
	rss_context uint32
	indir_size  uint32
	key_size    uint32
	hfunc       byte
	rsvd8       [3]byte
	rsvd32      uint32
	//__u32 rss_config[0];
}
type ethtoolTcpip4Spec struct {
	ip4src be32
	ip4dst be32
	psrc   be16
	pdst   be16
	tos    uint8
}

type ethtoolTcpip6Spec struct {
	ip6src [4]be32
	ip6dst [4]be32
	psrc   be16
	pdst   be16
	tclass byte
}

type ethtoolAhEspip4Spec struct {
	ip4src be32
	ip4dst be32
	spi    be32
	tos    byte
}

type ethtoolAhEspip6Spec struct {
	ip6src [4]be32
	ip6dst [4]be32
	spi    be32
	tclass byte
}

type ethtoolUsrip4Spec struct {
	ip4src     be32
	ip4dst     be32
	l4_4_bytes be32
	tos        byte
	ip_ver     byte
	proto      byte
}

type ethtoolUsrip6Spec struct {
	ip6src     [4]be32
	ip6dst     [4]be32
	l4_4_bytes be32
	tclass     byte
	l4_proto   byte
}

type ethhdr struct {
	h_dest   [ETH_ALEN]byte // destination eth addr
	h_source [ETH_ALEN]byte // source ether addr
	h_proto  be16           // packet type ID field
}

type ethtoolFlowUnion struct {
	// struct ethtool_tcpip4_spec		tcp_ip4_spec;
	// struct ethtool_tcpip4_spec		udp_ip4_spec;
	// struct ethtool_tcpip4_spec		sctp_ip4_spec;
	// struct ethtool_ah_espip4_spec		ah_ip4_spec;
	// struct ethtool_ah_espip4_spec		esp_ip4_spec;
	// struct ethtool_usrip4_spec		usr_ip4_spec;
	// struct ethtool_tcpip6_spec		tcp_ip6_spec;
	// struct ethtool_tcpip6_spec		udp_ip6_spec;
	// struct ethtool_tcpip6_spec		sctp_ip6_spec;
	// struct ethtool_ah_espip6_spec		ah_ip6_spec;
	// struct ethtool_ah_espip6_spec		esp_ip6_spec;
	// struct ethtool_usrip6_spec		usr_ip6_spec;
	// struct ethhdr				ether_spec;
	hdata [52]byte
}

func (u *ethtoolFlowUnion) tcpIp4Spec() *ethtoolTcpip4Spec {
	return (*ethtoolTcpip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) udpIp4Spec() *ethtoolTcpip4Spec {
	return (*ethtoolTcpip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) sctpIp4Spec() *ethtoolTcpip4Spec {
	return (*ethtoolTcpip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) ahIp4Spec() *ethtoolAhEspip4Spec {
	return (*ethtoolAhEspip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) espIp4Spec() *ethtoolAhEspip4Spec {
	return (*ethtoolAhEspip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) usrIp4Spec() *ethtoolUsrip4Spec {
	return (*ethtoolUsrip4Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) tcpIp6Spec() *ethtoolTcpip6Spec {
	return (*ethtoolTcpip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) udpIp6Spec() *ethtoolTcpip6Spec {
	return (*ethtoolTcpip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) sctpIp6Spec() *ethtoolTcpip6Spec {
	return (*ethtoolTcpip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) ahIp6Spec() *ethtoolAhEspip6Spec {
	return (*ethtoolAhEspip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) espIp6Spec() *ethtoolAhEspip6Spec {
	return (*ethtoolAhEspip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) usrIp6Spec() *ethtoolUsrip6Spec {
	return (*ethtoolUsrip6Spec)(unsafe.Pointer(&u.hdata[0]))
}

func (u *ethtoolFlowUnion) etherSpec() *ethhdr {
	return (*ethhdr)(unsafe.Pointer(&u.hdata[0]))
}

type ethtoolFlowExt struct {
	padding    [2]byte
	h_dest     [ETH_ALEN]byte
	vlan_etype be16
	vlan_tci   be16
	data       [2]be32
}

type ethtoolRxFlowSpec struct {
	flow_type   uint32
	h_u         ethtoolFlowUnion
	h_ext       ethtoolFlowExt
	m_u         ethtoolFlowUnion
	m_ext       ethtoolFlowExt
	ring_cookie uint64
	location    uint32
}

type ethtoolRxnfc struct {
	cmd                     uint32
	flow_type               uint32
	data                    uint64
	fs                      ethtoolRxFlowSpec
	rule_cnt_or_rss_context uint32
}

type ethtoolRxfhIndir struct {
	cmd  uint32
	size uint32
	//__u32 ring_index[0];
}

type Ethtool struct {
	fd int
}

// Convert zero-terminated array of chars (string in C) to a Go string.
func goString(s []byte) string {
	strEnd := bytes.IndexByte(s, 0)
	if strEnd == -1 {
		return string(s[:])
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

// ModuleEeprom returns Eeprom information of the given interface name.
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

func (e *Ethtool) getDriverInfo(intf string) (*ethtoolDrvInfo, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return nil, err
	}

	return &drvinfo, nil
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

func isFeatureBitSet(blocks [MAX_FEATURE_BLOCKS]ethtoolGetFeaturesBlock, index uint) bool {
	return (blocks)[index/32].active&(1<<(index%32)) != 0
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

// FeatureNames shows supported features by their name.
func (e *Ethtool) FeatureNames(intf string) (map[string]uint, error) {
	return e.getStringSet(intf, ETH_SS_FEATURES, 0)
}

func (e *Ethtool) getStringSet(intf string, ss stringSet, drvinfoOffset uintptr) (map[string]uint, error) {
	ssetInfo := ethtoolSsetInfo{
		cmd:       ETHTOOL_GSSET_INFO,
		sset_mask: 1 << ss,
	}

	var length uint32

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ssetInfo))); err == nil {
		if ssetInfo.sset_mask != 0 {
			length = uint32(ssetInfo.data)
		}
	} else if err == syscall.EOPNOTSUPP && drvinfoOffset != 0 {
		drvinfo, err := e.getDriverInfo(intf)
		if err != nil {
			return nil, err
		}

		length = *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(drvinfo)) + drvinfoOffset))
	} else {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	} else if length > MAX_GSTRINGS {
		return nil, fmt.Errorf("ethtool currently doesn't support more than %d entries, received %d", MAX_GSTRINGS, length)
	}

	gstrings := ethtoolGStrings{
		cmd:        ETHTOOL_GSTRINGS,
		string_set: uint32(ss),
		len:        length,
		data:       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&gstrings))); err != nil {
		return nil, err
	}

	var result = make(map[string]uint)
	for i := 0; i != int(length); i++ {
		b := gstrings.data[i*ETH_GSTRING_LEN : i*ETH_GSTRING_LEN+ETH_GSTRING_LEN]
		key := goString(b)
		if key != "" {
			result[key] = uint(i)
		}
	}

	return result, nil
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

	var result = make(map[string]bool, length)
	for key, index := range names {
		result[key] = isFeatureBitSet(features.blocks, index)
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

// Get state of a link.
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
func (e *Ethtool) Stats(intf string) (map[string]uint64, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return nil, err
	}

	if drvinfo.n_stats*ETH_GSTRING_LEN > MAX_GSTRINGS*ETH_GSTRING_LEN {
		return nil, fmt.Errorf("ethtool currently doesn't support more than %d entries, received %d", MAX_GSTRINGS, drvinfo.n_stats)
	}

	gstrings := ethtoolGStrings{
		cmd:        ETHTOOL_GSTRINGS,
		string_set: uint32(ETH_SS_STATS),
		len:        drvinfo.n_stats,
		data:       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&gstrings))); err != nil {
		return nil, err
	}

	stats := ethtoolStats{
		cmd:     ETHTOOL_GSTATS,
		n_stats: drvinfo.n_stats,
		data:    [MAX_GSTRINGS]uint64{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&stats))); err != nil {
		return nil, err
	}

	var result = make(map[string]uint64)
	for i := 0; i != int(drvinfo.n_stats); i++ {
		b := gstrings.data[i*ETH_GSTRING_LEN : i*ETH_GSTRING_LEN+ETH_GSTRING_LEN]
		strEnd := strings.Index(string(b), "\x00")
		if strEnd == -1 {
			strEnd = ETH_GSTRING_LEN
		}
		key := string(b[:strEnd])
		if len(key) != 0 {
			result[key] = stats.data[i]
		}
	}

	return result, nil
}

type IndirectTable []uint32

func (t IndirectTable) String() string {
	var b strings.Builder

	for i, n := range t {
		if i%8 == 0 {
			fmt.Fprintf(&b, "%5d: ", i)
		}

		fmt.Fprintf(&b, " %5d", n)

		if i%8 == 7 || i == len(t)-1 {
			fmt.Fprintln(&b, "")
		}
	}

	return b.String()
}

type FlowHash struct {
	RingCount int
	Key       []byte
	Funcs     map[string]bool
	Table     IndirectTable
}

type flowHashConfig struct {
	rss_context uint32
}

type FlowHashOption func(*flowHashConfig)

func WithRSSContext(context uint32) FlowHashOption {
	return func(c *flowHashConfig) {
		c.rss_context = context
	}
}

// GetFlowHash get rx flow hash indirection table and/or RSS hash key
func (e *Ethtool) GetFlowHash(intf string, opts ...FlowHashOption) (*FlowHash, error) {
	var cfg flowHashConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	ringCount := ethtoolRxnfc{cmd: ETHTOOL_GRXRINGS}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ringCount))); err != nil {
		return nil, fmt.Errorf("get RX ring count, %w", err)
	}

	rssHead := ethtoolRxfh{cmd: ETHTOOL_GRSSH, rss_context: cfg.rss_context}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&rssHead))); err != nil {
		if err == syscall.EOPNOTSUPP && cfg.rss_context != 0 {
			table, err := e.getFlowHashIndirectTable(intf)
			if err != nil {
				return nil, err
			}
			return &FlowHash{RingCount: int(ringCount.data), Table: table}, nil
		}

		return nil, fmt.Errorf("get RX flow hash indir size and/or key size, %w", err)
	}

	sz := unsafe.Sizeof(ethtoolRxfh{}) + uintptr(rssHead.indir_size)*unsafe.Sizeof(uint32(0)) + uintptr(rssHead.key_size)
	rss := (*ethtoolRxfh)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(rss))

	rss.cmd = ETHTOOL_GRSSH
	rss.rss_context = cfg.rss_context
	rss.indir_size = rssHead.indir_size
	rss.key_size = rssHead.key_size

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(rss))); err != nil {
		return nil, fmt.Errorf("get RX flow hash configuration, %w", err)
	}

	n := uintptr(rss.indir_size)
	if unsafe.Sizeof(ethtoolRxfh{})+n*unsafe.Sizeof(uint32(0)) > sz {
		return nil, fmt.Errorf("get RX flow indirect table, %w", syscall.ERANGE)
	}

	p := unsafe.Pointer(uintptr(unsafe.Pointer(rss)) + unsafe.Sizeof(ethtoolRxfh{}))
	table := make([]uint32, n)
	copy(table, (*[1 << 24]uint32)(p)[:n])

	var key []byte

	if rss.key_size > 0 {
		off := unsafe.Sizeof(ethtoolRxfh{}) + uintptr(rss.indir_size)*unsafe.Sizeof(uint32(0))

		if off+uintptr(rss.key_size) > sz {
			return nil, syscall.ERANGE
		} else {
			key = C.GoBytes(unsafe.Pointer((uintptr(unsafe.Pointer(rss)) + off)), C.int(rss.key_size))
		}
	}

	var funcs map[string]bool

	if rss.hfunc == 0 {
		return nil, syscall.ENOTSUP
	}

	hfuncs, err := e.getStringSet(intf, ETH_SS_RSS_HASH_FUNCS, 0)
	if err != nil {
		return nil, fmt.Errorf("get hash functions names, %w", err)
	}

	funcs = make(map[string]bool, bits.OnesCount8(rss.hfunc))
	for name, i := range hfuncs {
		funcs[name] = (rss.hfunc & (1 << i)) != 0
	}

	return &FlowHash{
		RingCount: int(ringCount.data),
		Key:       key,
		Funcs:     funcs,
		Table:     table,
	}, nil
}

func (e *Ethtool) getFlowHashIndirectTable(intf string) ([]uint32, error) {
	indirHead := ethtoolRxfhIndir{cmd: ETHTOOL_GRXFHINDIR}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indirHead))); err != nil {
		return nil, fmt.Errorf("get RX flow hash indirection table size, %w", err)
	}

	sz := unsafe.Sizeof(ethtoolRxfhIndir{}) + uintptr(indirHead.size)*unsafe.Sizeof(uint32(0))
	indir := (*ethtoolRxfhIndir)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(indir))

	indir.cmd = ETHTOOL_GRXFHINDIR
	indir.size = indirHead.size

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indirHead))); err != nil {
		return nil, fmt.Errorf("get RX flow hash indirection table, %w", err)
	}

	if indir.size == 0 {
		return nil, syscall.ENOTSUP
	}

	n := uintptr(indir.size)
	if unsafe.Sizeof(ethtoolRxfhIndir{})+n*unsafe.Sizeof(uint32(0)) > sz {
		return nil, fmt.Errorf("get RX flow indirect table, %w", syscall.ERANGE)
	}

	p := unsafe.Pointer(uintptr(unsafe.Pointer(indir)) + unsafe.Sizeof(ethtoolRxfhIndir{}))
	table := make([]uint32, n)
	copy(table, (*[1 << 24]uint32)(p)[:n])

	return table, nil
}

// Close closes the ethool handler
func (e *Ethtool) Close() {
	unix.Close(e.fd)
}

// NewEthtool returns a new ethtool handler
func NewEthtool() (*Ethtool, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_IP)
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
