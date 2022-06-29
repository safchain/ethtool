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
	"errors"
	"fmt"
	"math/bits"
	"sort"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/safchain/ethtool/flowhash"
)

// #include <string.h>
// #include <stdlib.h>
import "C"

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
)

const (
	// MAX_GSTRINGS maximum number of stats entries that ethtool can retrieve currently.
	MAX_GSTRINGS       = 32768
	MAX_FEATURE_BLOCKS = (MAX_GSTRINGS + 32 - 1) / 32
	EEPROM_LEN         = 640
	PERMADDR_LEN       = 32
	ETH_ALEN           = 6
)

/* Duplex, half or full. */
const (
	DUPLEX_HALF    = 0x00
	DUPLEX_FULL    = 0x01
	DUPLEX_UNKNOWN = 0xff
)

var duplexNames = map[uint8]string{
	DUPLEX_HALF: "half",
	DUPLEX_FULL: "full",
}

func DuplexName(v uint8) string {
	s, ok := duplexNames[v]
	if ok {
		return s
	}
	return fmt.Sprintln("Unknown", v)
}

/* Which connector port. */
const (
	PORT_TP    = 0x00
	PORT_AUI   = 0x01
	PORT_MII   = 0x02
	PORT_FIBRE = 0x03
	PORT_BNC   = 0x04
	PORT_DA    = 0x05
	PORT_NONE  = 0xef
	PORT_OTHER = 0xff
)

var portTypeNames = map[uint8]string{
	PORT_TP:    "TP",
	PORT_AUI:   "AUI",
	PORT_MII:   "MII",
	PORT_FIBRE: "FIBRE",
	PORT_BNC:   "BNC",
	PORT_DA:    "Direct Attach Copper",
	PORT_NONE:  "None",
	PORT_OTHER: "Other",
}

func PortTypeName(v uint8) string {
	s, ok := portTypeNames[v]
	if ok {
		return s
	}
	return fmt.Sprintln("Unknown", v)
}

/* Which transceiver to use. */
const (
	XCVR_INTERNAL = 0x00 /* PHY and MAC are in the same package */
	XCVR_EXTERNAL = 0x01 /* PHY and MAC are in different packages */
	XCVR_DUMMY1   = 0x02
	XCVR_DUMMY2   = 0x03
	XCVR_DUMMY3   = 0x04
)

func TransceiverName(v uint8) string {
	switch v {
	case XCVR_INTERNAL:
		return "internal"
	case XCVR_EXTERNAL:
		return "external"
	default:
		return fmt.Sprintln("unknown", v)
	}
}

var linkSpeedNames = map[uint32]string{
	unix.ETHTOOL_LINK_MODE_10baseT_Half_BIT:               "10baseT/Half",
	unix.ETHTOOL_LINK_MODE_10baseT_Full_BIT:               "10baseT/Full",
	unix.ETHTOOL_LINK_MODE_100baseT_Half_BIT:              "100baseT/Half",
	unix.ETHTOOL_LINK_MODE_100baseT_Full_BIT:              "100baseT/Full",
	unix.ETHTOOL_LINK_MODE_1000baseT_Half_BIT:             "1000baseT/Half",
	unix.ETHTOOL_LINK_MODE_1000baseT_Full_BIT:             "1000baseT/Full",
	unix.ETHTOOL_LINK_MODE_10000baseT_Full_BIT:            "10000baseT/Full",
	unix.ETHTOOL_LINK_MODE_2500baseX_Full_BIT:             "2500baseX/Full",
	unix.ETHTOOL_LINK_MODE_1000baseKX_Full_BIT:            "1000baseKX/Full",
	unix.ETHTOOL_LINK_MODE_10000baseKX4_Full_BIT:          "10000baseKX4/Full",
	unix.ETHTOOL_LINK_MODE_10000baseKR_Full_BIT:           "10000baseKR/Full",
	unix.ETHTOOL_LINK_MODE_10000baseR_FEC_BIT:             "10000baseR_FEC",
	unix.ETHTOOL_LINK_MODE_20000baseMLD2_Full_BIT:         "20000baseMLD2/Full",
	unix.ETHTOOL_LINK_MODE_20000baseKR2_Full_BIT:          "20000baseKR2/Full",
	unix.ETHTOOL_LINK_MODE_40000baseKR4_Full_BIT:          "40000baseKR4/Full",
	unix.ETHTOOL_LINK_MODE_40000baseCR4_Full_BIT:          "40000baseCR4/Full",
	unix.ETHTOOL_LINK_MODE_40000baseSR4_Full_BIT:          "40000baseSR4/Full",
	unix.ETHTOOL_LINK_MODE_40000baseLR4_Full_BIT:          "40000baseLR4/Full",
	unix.ETHTOOL_LINK_MODE_56000baseKR4_Full_BIT:          "56000baseKR4/Full",
	unix.ETHTOOL_LINK_MODE_56000baseCR4_Full_BIT:          "56000baseCR4/Full",
	unix.ETHTOOL_LINK_MODE_56000baseSR4_Full_BIT:          "56000baseSR4/Full",
	unix.ETHTOOL_LINK_MODE_56000baseLR4_Full_BIT:          "56000baseLR4/Full",
	unix.ETHTOOL_LINK_MODE_25000baseCR_Full_BIT:           "25000baseCR/Full",
	unix.ETHTOOL_LINK_MODE_25000baseKR_Full_BIT:           "25000baseKR/Full",
	unix.ETHTOOL_LINK_MODE_25000baseSR_Full_BIT:           "25000baseSR/Full",
	unix.ETHTOOL_LINK_MODE_50000baseCR2_Full_BIT:          "50000baseCR2/Full",
	unix.ETHTOOL_LINK_MODE_50000baseKR2_Full_BIT:          "50000baseKR2/Full",
	unix.ETHTOOL_LINK_MODE_100000baseKR4_Full_BIT:         "100000baseKR4/Full",
	unix.ETHTOOL_LINK_MODE_100000baseSR4_Full_BIT:         "100000baseSR4/Full",
	unix.ETHTOOL_LINK_MODE_100000baseCR4_Full_BIT:         "100000baseCR4/Full",
	unix.ETHTOOL_LINK_MODE_100000baseLR4_ER4_Full_BIT:     "100000baseLR4_ER4/Full",
	unix.ETHTOOL_LINK_MODE_50000baseSR2_Full_BIT:          "50000baseSR2/Full",
	unix.ETHTOOL_LINK_MODE_1000baseX_Full_BIT:             "1000baseX/Full",
	unix.ETHTOOL_LINK_MODE_10000baseCR_Full_BIT:           "10000baseCR/Full",
	unix.ETHTOOL_LINK_MODE_10000baseSR_Full_BIT:           "10000baseSR/Full",
	unix.ETHTOOL_LINK_MODE_10000baseLR_Full_BIT:           "10000baseLR/Full",
	unix.ETHTOOL_LINK_MODE_10000baseLRM_Full_BIT:          "10000baseLRM/Full",
	unix.ETHTOOL_LINK_MODE_10000baseER_Full_BIT:           "10000baseER/Full",
	unix.ETHTOOL_LINK_MODE_2500baseT_Full_BIT:             "2500baseT/Full",
	unix.ETHTOOL_LINK_MODE_5000baseT_Full_BIT:             "5000baseT/Full",
	unix.ETHTOOL_LINK_MODE_50000baseKR_Full_BIT:           "50000baseKR/Full",
	unix.ETHTOOL_LINK_MODE_50000baseSR_Full_BIT:           "50000baseSR/Full",
	unix.ETHTOOL_LINK_MODE_50000baseCR_Full_BIT:           "50000baseCR/Full",
	unix.ETHTOOL_LINK_MODE_50000baseLR_ER_FR_Full_BIT:     "50000baseLR_ER_FR/Full",
	unix.ETHTOOL_LINK_MODE_50000baseDR_Full_BIT:           "50000baseDR/Full",
	unix.ETHTOOL_LINK_MODE_100000baseKR2_Full_BIT:         "100000baseKR2/Full",
	unix.ETHTOOL_LINK_MODE_100000baseSR2_Full_BIT:         "100000baseSR2/Full",
	unix.ETHTOOL_LINK_MODE_100000baseCR2_Full_BIT:         "100000baseCR2/Full",
	unix.ETHTOOL_LINK_MODE_100000baseLR2_ER2_FR2_Full_BIT: "100000baseLR2_ER2_FR2/Full",
	unix.ETHTOOL_LINK_MODE_100000baseDR2_Full_BIT:         "100000baseDR2/Full",
	unix.ETHTOOL_LINK_MODE_200000baseKR4_Full_BIT:         "200000baseKR4/Full",
	unix.ETHTOOL_LINK_MODE_200000baseSR4_Full_BIT:         "200000baseSR4/Full",
	unix.ETHTOOL_LINK_MODE_200000baseLR4_ER4_FR4_Full_BIT: "200000baseLR4_ER4_FR4/Full",
	unix.ETHTOOL_LINK_MODE_200000baseDR4_Full_BIT:         "200000baseDR4/Full",
	unix.ETHTOOL_LINK_MODE_200000baseCR4_Full_BIT:         "200000baseCR4/Full",
	unix.ETHTOOL_LINK_MODE_100baseT1_Full_BIT:             "100baseT1/Full",
	unix.ETHTOOL_LINK_MODE_1000baseT1_Full_BIT:            "1000baseT1/Full",
	unix.ETHTOOL_LINK_MODE_400000baseKR8_Full_BIT:         "400000baseKR8/Full",
	unix.ETHTOOL_LINK_MODE_400000baseSR8_Full_BIT:         "400000baseSR8/Full",
	unix.ETHTOOL_LINK_MODE_400000baseLR8_ER8_FR8_Full_BIT: "400000baseLR8_ER8_FR8/Full",
	unix.ETHTOOL_LINK_MODE_400000baseDR8_Full_BIT:         "400000baseDR8/Full",
	unix.ETHTOOL_LINK_MODE_400000baseCR8_Full_BIT:         "400000baseCR8/Full",
	unix.ETHTOOL_LINK_MODE_100000baseKR_Full_BIT:          "100000baseKR/Full",
	unix.ETHTOOL_LINK_MODE_100000baseSR_Full_BIT:          "100000baseSR/Full",
	unix.ETHTOOL_LINK_MODE_100000baseLR_ER_FR_Full_BIT:    "100000baseLR_ER_FR/Full",
	unix.ETHTOOL_LINK_MODE_100000baseDR_Full_BIT:          "100000baseDR/Full",
	unix.ETHTOOL_LINK_MODE_100000baseCR_Full_BIT:          "100000baseCR/Full",
	unix.ETHTOOL_LINK_MODE_200000baseKR2_Full_BIT:         "200000baseKR2/Full",
	unix.ETHTOOL_LINK_MODE_200000baseSR2_Full_BIT:         "200000baseSR2/Full",
	unix.ETHTOOL_LINK_MODE_200000baseLR2_ER2_FR2_Full_BIT: "200000baseLR2_ER2_FR2/Full",
	unix.ETHTOOL_LINK_MODE_200000baseDR2_Full_BIT:         "200000baseDR2/Full",
	unix.ETHTOOL_LINK_MODE_200000baseCR2_Full_BIT:         "200000baseCR2/Full",
	unix.ETHTOOL_LINK_MODE_400000baseKR4_Full_BIT:         "400000baseKR4/Full",
	unix.ETHTOOL_LINK_MODE_400000baseSR4_Full_BIT:         "400000baseSR4/Full",
	unix.ETHTOOL_LINK_MODE_400000baseLR4_ER4_FR4_Full_BIT: "400000baseLR4_ER4_FR4/Full",
	unix.ETHTOOL_LINK_MODE_400000baseDR4_Full_BIT:         "400000baseDR4/Full",
	unix.ETHTOOL_LINK_MODE_400000baseCR4_Full_BIT:         "400000baseCR4/Full",
	unix.ETHTOOL_LINK_MODE_100baseFX_Half_BIT:             "100baseFX/Half",
	unix.ETHTOOL_LINK_MODE_100baseFX_Full_BIT:             "100baseFX/Full",
}

func LinkModeTestBit(v []uint32, bit int) bool {
	if bit/8 >= len(v) {
		return false
	}
	return v[bit/32]&(1<<(bit%32)) != 0
}

func LinkPortNames(v []uint32) (names []string) {
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_TP_BIT) {
		names = append(names, "TP")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_AUI_BIT) {
		names = append(names, "AUI")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_BNC_BIT) {
		names = append(names, "BNC")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_MII_BIT) {
		names = append(names, "MII")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_FIBRE_BIT) {
		names = append(names, "FIBRE")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_Backplane_BIT) {
		names = append(names, "Backplane")
	}

	return
}

func LinkECCModeNames(v []uint32) (names []string) {
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_FEC_NONE_BIT) {
		names = append(names, "None")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_FEC_BASER_BIT) {
		names = append(names, "BaseR")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_FEC_RS_BIT) {
		names = append(names, "RS")
	}
	if LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_FEC_LLRS_BIT) {
		names = append(names, "LLRS")
	}

	return
}

func LinkSpeedNames(v []uint32) (names []string) {
	var s []int
	for speed := range linkSpeedNames {
		s = append(s, int(speed))
	}

	sort.Ints(s)

	for _, n := range s {
		if LinkModeTestBit(v, n) {
			names = append(names, linkSpeedNames[uint32(n)])
		}
	}

	return
}

// MDI or MDI-X status/control - if MDI/MDI_X/AUTO is set then the driver is required to renegotiate link

const (
	ETH_TP_MDI_INVALID = 0x00 // status: unknown; control: unsupported
	ETH_TP_MDI         = 0x01 // status: MDI;     control: force MDI
	ETH_TP_MDI_X       = 0x02 // status: MDI-X;   control: force MDI-X
	ETH_TP_MDI_AUTO    = 0x03 //                  control: auto-select
)

const (
	NETIF_MSG_DRV       = 0x0001
	NETIF_MSG_PROBE     = 0x0002
	NETIF_MSG_LINK      = 0x0004
	NETIF_MSG_TIMER     = 0x0008
	NETIF_MSG_IFDOWN    = 0x0010
	NETIF_MSG_IFUP      = 0x0020
	NETIF_MSG_RX_ERR    = 0x0040
	NETIF_MSG_TX_ERR    = 0x0080
	NETIF_MSG_TX_QUEUED = 0x0100
	NETIF_MSG_INTR      = 0x0200
	NETIF_MSG_TX_DONE   = 0x0400
	NETIF_MSG_RX_STATUS = 0x0800
	NETIF_MSG_PKTDATA   = 0x1000
	NETIF_MSG_HW        = 0x2000
	NETIF_MSG_WOL       = 0x4000
)

var msgLevelNames = map[uint32]string{
	NETIF_MSG_DRV:       "drv",
	NETIF_MSG_PROBE:     "probe",
	NETIF_MSG_LINK:      "link",
	NETIF_MSG_TIMER:     "timer",
	NETIF_MSG_IFDOWN:    "ifdown",
	NETIF_MSG_IFUP:      "ifup",
	NETIF_MSG_RX_ERR:    "rx_err",
	NETIF_MSG_TX_ERR:    "tx_err",
	NETIF_MSG_TX_QUEUED: "tx_queued",
	NETIF_MSG_INTR:      "intr",
	NETIF_MSG_TX_DONE:   "tx_done",
	NETIF_MSG_RX_STATUS: "rx_status",
	NETIF_MSG_PKTDATA:   "pktdata",
	NETIF_MSG_HW:        "hw",
	NETIF_MSG_WOL:       "wol",
}

func MsgLevelNames(v uint32) (names []string) {
	for k, n := range msgLevelNames {
		if v&k != 0 {
			names = append(names, n)
		}
	}

	return
}

type ifreq struct {
	ifr_name [unix.IFNAMSIZ]byte
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
	var name [unix.IFNAMSIZ]byte
	copy(name[:], []byte(intf))

	ifr := ifreq{
		ifr_name: name,
		ifr_data: data,
	}

	_, _, ep := unix.Syscall(unix.SYS_IOCTL, uintptr(e.fd), unix.SIOCETHTOOL, uintptr(unsafe.Pointer(&ifr)))
	if ep != 0 {
		return ep
	}

	return nil
}

func (e *Ethtool) getDriverInfo(intf string) (*ethtoolDrvInfo, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: unix.ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return nil, err
	}

	return &drvinfo, nil
}

func (e *Ethtool) getChannels(intf string) (Channels, error) {
	channels := Channels{
		Cmd: unix.ETHTOOL_GCHANNELS,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&channels))); err != nil {
		return Channels{}, err
	}

	return channels, nil
}

func (e *Ethtool) setChannels(intf string, channels Channels) (Channels, error) {
	channels.Cmd = unix.ETHTOOL_SCHANNELS

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&channels))); err != nil {
		return Channels{}, err
	}

	return channels, nil
}

func (e *Ethtool) getCoalesce(intf string) (Coalesce, error) {
	coalesce := Coalesce{
		Cmd: unix.ETHTOOL_GCOALESCE,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&coalesce))); err != nil {
		return Coalesce{}, err
	}

	return coalesce, nil
}

func (e *Ethtool) getPermAddr(intf string) (ethtoolPermAddr, error) {
	permAddr := ethtoolPermAddr{
		cmd:  unix.ETHTOOL_GPERMADDR,
		size: PERMADDR_LEN,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&permAddr))); err != nil {
		return ethtoolPermAddr{}, err
	}

	return permAddr, nil
}

func (e *Ethtool) getModuleEeprom(intf string) (ethtoolEeprom, ethtoolModInfo, error) {
	modInfo := ethtoolModInfo{
		cmd: unix.ETHTOOL_GMODULEINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&modInfo))); err != nil {
		return ethtoolEeprom{}, ethtoolModInfo{}, err
	}

	eeprom := ethtoolEeprom{
		cmd:    unix.ETHTOOL_GMODULEEEPROM,
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
func (e *Ethtool) FeatureNames(intf string) (StringSet, error) {
	return e.getStringSet(intf, ETH_SS_FEATURES, 0)
}

type StringSet map[string]uint

func (e *Ethtool) getStringSet(intf string, ss stringSet, drvinfoOffset uintptr) (StringSet, error) {
	ssetInfo := ethtoolSsetInfo{
		cmd:       unix.ETHTOOL_GSSET_INFO,
		sset_mask: 1 << ss,
	}

	var length uint32

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ssetInfo))); err == nil {
		if ssetInfo.sset_mask != 0 {
			length = uint32(ssetInfo.data)
		}
	} else if errors.Is(err, syscall.EOPNOTSUPP) && drvinfoOffset != 0 {
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
		cmd:        unix.ETHTOOL_GSTRINGS,
		string_set: uint32(ss),
		len:        length,
		data:       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&gstrings))); err != nil {
		return nil, err
	}

	var result = make(StringSet)
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
		cmd:  unix.ETHTOOL_GFEATURES,
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
		cmd:  unix.ETHTOOL_SFEATURES,
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
		cmd: unix.ETHTOOL_GLINK,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&x))); err != nil {
		return 0, err
	}

	return x.data, nil
}

// Stats retrieves stats of the given interface name.
func (e *Ethtool) Stats(intf string) (map[string]uint64, error) {
	drvinfo := ethtoolDrvInfo{
		cmd: unix.ETHTOOL_GDRVINFO,
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&drvinfo))); err != nil {
		return nil, err
	}

	if drvinfo.n_stats*ETH_GSTRING_LEN > MAX_GSTRINGS*ETH_GSTRING_LEN {
		return nil, fmt.Errorf("ethtool currently doesn't support more than %d entries, received %d", MAX_GSTRINGS, drvinfo.n_stats)
	}

	gstrings := ethtoolGStrings{
		cmd:        unix.ETHTOOL_GSTRINGS,
		string_set: uint32(ETH_SS_STATS),
		len:        drvinfo.n_stats,
		data:       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte{},
	}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&gstrings))); err != nil {
		return nil, err
	}

	stats := ethtoolStats{
		cmd:     unix.ETHTOOL_GSTATS,
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

// GetFlowHash get rx flow hash indirection table and/or RSS hash key
func (e *Ethtool) GetFlowHash(intf string, opts ...flowhash.Option) (*flowhash.FlowHash, error) {
	o := flowhash.NewConfig(opts)

	ringCount := ethtoolRxnfc{cmd: unix.ETHTOOL_GRXRINGS}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&ringCount))); err != nil {
		return nil, fmt.Errorf("get RX ring count, %w", err)
	}

	rssHead := ethtoolRxfh{cmd: unix.ETHTOOL_GRSSH, rss_context: uint32(o.Context)}

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&rssHead))); err != nil {
		if errors.Is(err, syscall.EOPNOTSUPP) && o.Context != 0 {
			table, err := e.getFlowHashIndirectTable(intf)
			if err != nil {
				return nil, err
			}
			return &flowhash.FlowHash{RingCount: int(ringCount.data), Table: table}, nil
		}

		return nil, fmt.Errorf("get RX flow hash indir size and/or key size, %w", err)
	}

	table := flowhash.NewIndirectTable(int(rssHead.indir_size))

	sz := sizeofEthtoolRxfh + table.Size() + uintptr(rssHead.key_size)
	rss := (*ethtoolRxfh)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(rss))

	rss.cmd = unix.ETHTOOL_GRSSH
	rss.rss_context = uint32(o.Context)
	rss.indir_size = rssHead.indir_size
	rss.key_size = rssHead.key_size

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(rss))); err != nil {
		return nil, fmt.Errorf("get RX flow hash configuration, %w", err)
	}

	if sizeofEthtoolRxfh+table.Size() > sz {
		return nil, fmt.Errorf("get RX flow indirect table, %w", syscall.ERANGE)
	}

	p := unsafe.Pointer(uintptr(unsafe.Pointer(rss)) + sizeofEthtoolRxfh)
	copy(table, flowhash.UnsafeRawIndirectTable(p, int(rss.indir_size)))

	var key []byte

	if rss.key_size > 0 {
		off := sizeofEthtoolRxfh + table.Size()

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

	return &flowhash.FlowHash{
		RingCount: int(ringCount.data),
		Key:       key,
		Funcs:     funcs,
		Table:     table,
	}, nil
}

func (e *Ethtool) getFlowHashIndirectTable(intf string) (table flowhash.IndirectTable, err error) {
	indirHead := ethtoolRxfhIndir{cmd: unix.ETHTOOL_GRXFHINDIR}

	if err = e.ioctl(intf, uintptr(unsafe.Pointer(&indirHead))); err != nil {
		err = fmt.Errorf("get RX flow hash indirection table size, %w", err)
		return
	}

	table = flowhash.NewIndirectTable(int(indirHead.size))

	sz := sizeofEthtoolRxfhIndir + table.Size()
	indir := (*ethtoolRxfhIndir)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(indir))

	indir.cmd = unix.ETHTOOL_GRXFHINDIR
	indir.size = indirHead.size

	if err = e.ioctl(intf, uintptr(unsafe.Pointer(&indirHead))); err != nil {
		err = fmt.Errorf("get RX flow hash indirection table, %w", err)
		return
	}

	if indir.size == 0 {
		err = syscall.ENOTSUP
		return
	}

	if sizeofEthtoolRxfhIndir+table.Size() > sz {
		err = fmt.Errorf("get RX flow indirect table, %w", syscall.ERANGE)
		return
	}

	p := unsafe.Pointer(uintptr(unsafe.Pointer(indir)) + sizeofEthtoolRxfhIndir)
	copy(table, flowhash.UnsafeRawIndirectTable(p, int(indir.size)))

	return
}

const sizeofEthtoolRxfh = unsafe.Sizeof(ethtoolRxfh{})

func (e *Ethtool) SetFlowHash(intf string, opts ...flowhash.SetOption) (ctxt flowhash.RSSContext, err error) {
	c := flowhash.NewSetConfig(opts)

	ringCount := ethtoolRxnfc{cmd: unix.ETHTOOL_GRXRINGS}

	if err = e.ioctl(intf, uintptr(unsafe.Pointer(&ringCount))); err != nil {
		err = fmt.Errorf("get RX ring count, %w", err)
		return
	}

	rssHead := ethtoolRxfh{cmd: unix.ETHTOOL_GRSSH}
	if err = e.ioctl(intf, uintptr(unsafe.Pointer(&rssHead))); err != nil {
		if errors.Is(err, syscall.EOPNOTSUPP) && len(c.Key) == 0 && len(c.Func) == 0 && c.Action.(*flowhash.Delete) == nil {
			return 0, e.setFlowHashIndirect(intf, c)
		}

		err = fmt.Errorf("get RX flow hash indir size and key size, %w", err)
		return
	}

	var indirBytes uintptr

	if c.Action.(*flowhash.Equal) != nil || c.Action.(*flowhash.Weight) != nil {
		indirBytes = flowhash.IndirectTableSize(rssHead.indir_size)
	}

	var hfunc byte

	if rssHead.hfunc != 0 && len(c.Func) > 0 {
		var funcs StringSet
		funcs, err = e.getStringSet(intf, ETH_SS_RSS_HASH_FUNCS, 0)
		if err != nil {
			err = fmt.Errorf("get hash functions names, %w", err)
			return
		}

		if v, exists := funcs[c.Func]; exists {
			hfunc = 1 << v
		}

		if hfunc == 0 {
			err = fmt.Errorf("unknown hash function `%s`", c.Func)
			return
		}
	}

	sz := sizeofEthtoolRxfh + indirBytes + uintptr(rssHead.key_size)
	rss := (*ethtoolRxfh)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(rss))

	rss.cmd = unix.ETHTOOL_SRSSH
	rss.rss_context = uint32(c.Context)
	rss.hfunc = hfunc

	if c.Action.(*flowhash.Delete) == nil {
		rss.key_size = rssHead.key_size

		ptr := unsafe.Pointer(uintptr(unsafe.Pointer(rss)) + sizeofEthtoolRxfh)
		table := flowhash.UnsafeRawIndirectTable(ptr, int(rssHead.indir_size))

		var sz int
		sz, err = c.Action.Fill(table)
		if err != nil {
			err = fmt.Errorf("fill RX flow hash indirection table, %w", err)
			return
		}

		rss.indir_size = uint32(sz)
	}

	if len(c.Key) > 0 {
		dst := unsafe.Pointer(uintptr(unsafe.Pointer(rss)) + sizeofEthtoolRxfh + indirBytes)
		src := C.CBytes(c.Key)
		defer C.free(src)
		C.memcpy(dst, src, C.ulong(len(c.Key)))
	}

	if err = e.ioctl(intf, uintptr(unsafe.Pointer(rss))); err != nil {
		err = fmt.Errorf("set RX flow hash configuration, %w", err)
		return
	}

	if c.Context.IsNew() {
		ctxt = flowhash.RSSContext(rss.rss_context)
	}
	return
}

const sizeofEthtoolRxfhIndir = unsafe.Sizeof(ethtoolRxfhIndir{})

func (e *Ethtool) setFlowHashIndirect(intf string, c *flowhash.SetConfig) error {
	indirHead := ethtoolRxfhIndir{cmd: unix.ETHTOOL_GRXFHINDIR}
	if err := e.ioctl(intf, uintptr(unsafe.Pointer(&indirHead))); err != nil {
		return fmt.Errorf("get RX flow hash indirection table size, %w", err)
	}

	sz := sizeofEthtoolRxfhIndir + flowhash.IndirectTableSize(indirHead.size)
	indir := (*ethtoolRxfhIndir)(C.calloc(1, C.ulong(sz)))
	defer C.free(unsafe.Pointer(indir))

	indir.cmd = unix.ETHTOOL_SRXFHINDIR
	indir.size = indirHead.size

	ptr := unsafe.Pointer(uintptr(unsafe.Pointer(indir)) + sizeofEthtoolRxfhIndir)
	table := flowhash.UnsafeRawIndirectTable(ptr, int(indirHead.size))

	n, err := c.Action.Fill(table)
	if err != nil {
		return fmt.Errorf("fill RX flow hash indirection table, %w", err)
	}

	indir.size = uint32(n)

	if err := e.ioctl(intf, uintptr(unsafe.Pointer(indir))); err != nil {
		return fmt.Errorf("set RX flow hash indirection table, %w", err)
	}

	return nil
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
