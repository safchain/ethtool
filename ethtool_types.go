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
	"time"
)

// Maximum size of an interface name
const (
	IFNAMSIZ = 16
)

// ioctl ethtool request
const (
	SIOCETHTOOL = 0x8946
)

// ethtool stats related constants.
const (
	ETH_GSTRING_LEN   = 32
	ETH_SS_STATS      = 1
	ETH_SS_PRIV_FLAGS = 2
	ETH_SS_FEATURES   = 4

	// CMD supported
	ETHTOOL_GSET     = 0x00000001 /* Get settings. */
	ETHTOOL_SSET     = 0x00000002 /* Set settings. */
	ETHTOOL_GWOL     = 0x00000005 /* Get wake-on-lan options. */
	ETHTOOL_SWOL     = 0x00000006 /* Set wake-on-lan options. */
	ETHTOOL_GDRVINFO = 0x00000003 /* Get driver info. */
	ETHTOOL_GMSGLVL  = 0x00000007 /* Get driver message level */
	ETHTOOL_SMSGLVL  = 0x00000008 /* Set driver msg level. */

	// Get link status for host, i.e. whether the interface *and* the
	// physical port (if there is one) are up (ethtool_value).
	ETHTOOL_GLINK            = 0x0000000a
	ETHTOOL_GCOALESCE        = 0x0000000e /* Get coalesce config */
	ETHTOOL_SCOALESCE        = 0x0000000f /* Set coalesce config */
	ETHTOOL_GRINGPARAM       = 0x00000010 /* Get ring parameters */
	ETHTOOL_SRINGPARAM       = 0x00000011 /* Set ring parameters. */
	ETHTOOL_GPAUSEPARAM      = 0x00000012 /* Get pause parameters */
	ETHTOOL_SPAUSEPARAM      = 0x00000013 /* Set pause parameters. */
	ETHTOOL_GSTRINGS         = 0x0000001b /* Get specified string set */
	ETHTOOL_PHYS_ID          = 0x0000001c /* Identify the NIC */
	ETHTOOL_GSTATS           = 0x0000001d /* Get NIC-specific statistics */
	ETHTOOL_GPERMADDR        = 0x00000020 /* Get permanent hardware address */
	ETHTOOL_GFLAGS           = 0x00000025 /* Get flags bitmap(ethtool_value) */
	ETHTOOL_GPFLAGS          = 0x00000027 /* Get driver-private flags bitmap */
	ETHTOOL_SPFLAGS          = 0x00000028 /* Set driver-private flags bitmap */
	ETHTOOL_GSSET_INFO       = 0x00000037 /* Get string set info */
	ETHTOOL_GFEATURES        = 0x0000003a /* Get device offload settings */
	ETHTOOL_SFEATURES        = 0x0000003b /* Change device offload settings */
	ETHTOOL_GCHANNELS        = 0x0000003c /* Get no of channels */
	ETHTOOL_SCHANNELS        = 0x0000003d /* Set no of channels */
	ETHTOOL_GET_TS_INFO      = 0x00000041 /* Get time stamping and PHC info */
	ETHTOOL_GMODULEINFO      = 0x00000042 /* Get plug-in module information */
	ETHTOOL_GMODULEEEPROM    = 0x00000043 /* Get plug-in module eeprom */
	ETHTOOL_GRXFHINDIR       = 0x00000038 /* Get RX flow hash indir'n table */
	ETHTOOL_SRXFHINDIR       = 0x00000039 /* Set RX flow hash indir'n table */
	ETH_RXFH_INDIR_NO_CHANGE = 0xFFFFFFFF

	// Speed and Duplex unknowns/constants (Manually defined based on <linux/ethtool.h>)
	SPEED_UNKNOWN  = 0xffffffff // ((__u32)-1) SPEED_UNKNOWN
	DUPLEX_HALF    = 0x00       // DUPLEX_HALF
	DUPLEX_FULL    = 0x01       // DUPLEX_FULL
	DUPLEX_UNKNOWN = 0xff       // DUPLEX_UNKNOWN

	// Port types (Manually defined based on <linux/ethtool.h>)
	PORT_TP    = 0x00 // PORT_TP
	PORT_AUI   = 0x01 // PORT_AUI
	PORT_MII   = 0x02 // PORT_MII
	PORT_FIBRE = 0x03 // PORT_FIBRE
	PORT_BNC   = 0x04 // PORT_BNC
	PORT_DA    = 0x05 // PORT_DA
	PORT_NONE  = 0xef // PORT_NONE
	PORT_OTHER = 0xff // PORT_OTHER

	// Autoneg settings (Manually defined based on <linux/ethtool.h>)
	AUTONEG_DISABLE = 0x00 // AUTONEG_DISABLE
	AUTONEG_ENABLE  = 0x01 // AUTONEG_ENABLE

	// MDIX states (Manually defined based on <linux/ethtool.h>)
	ETH_TP_MDI_INVALID = 0x00 // ETH_TP_MDI_INVALID
	ETH_TP_MDI         = 0x01 // ETH_TP_MDI
	ETH_TP_MDI_X       = 0x02 // ETH_TP_MDI_X
	ETH_TP_MDI_AUTO    = 0x03 // Control value ETH_TP_MDI_AUTO

	// Link mode mask bits count (Manually defined based on ethtool.h)
	ETHTOOL_LINK_MODE_MASK_NBITS = 92 // __ETHTOOL_LINK_MODE_MASK_NBITS

	// Calculate max nwords based on NBITS using the manually defined constant
	MAX_LINK_MODE_MASK_NWORDS = (ETHTOOL_LINK_MODE_MASK_NBITS + 31) / 32 // = 3
)

// MAX_GSTRINGS maximum number of stats entries that ethtool can
// retrieve currently.
const (
	MAX_GSTRINGS       = 32768
	MAX_FEATURE_BLOCKS = (MAX_GSTRINGS + 32 - 1) / 32
	EEPROM_LEN         = 640
	PERMADDR_LEN       = 32
)

// ethtool sset_info related constants
const (
	MAX_SSET_INFO = 64
)

const (
	DEFAULT_BLINK_DURATION = 60 * time.Second
)

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

// IdentityConf is an identity config for an interface
type IdentityConf struct {
	Cmd      uint32
	Duration uint32
}

// WoL options
const (
	WAKE_PHY         = 1 << 0
	WAKE_UCAST       = 1 << 1
	WAKE_MCAST       = 1 << 2
	WAKE_BCAST       = 1 << 3
	WAKE_ARP         = 1 << 4
	WAKE_MAGIC       = 1 << 5
	WAKE_MAGICSECURE = 1 << 6 // only meaningful if WAKE_MAGIC
)

var WoLMap = map[uint32]string{
	WAKE_PHY:         "p", // Wake on PHY activity
	WAKE_UCAST:       "u", // Wake on unicast messages
	WAKE_MCAST:       "m", // Wake on multicast messages
	WAKE_BCAST:       "b", // Wake on broadcast messages
	WAKE_ARP:         "a", // Wake on ARP
	WAKE_MAGIC:       "g", // Wake on MagicPacket™
	WAKE_MAGICSECURE: "s", // Enable SecureOn™ password for MagicPacket™
	// f Wake on filter(s)
	// d Disable (wake on  nothing). This option clears all previous options.
}

// WakeOnLan contains WoL config for an interface
type WakeOnLan struct {
	Cmd       uint32 // ETHTOOL_GWOL or ETHTOOL_SWOL
	Supported uint32 // r/o bitmask of WAKE_* flags for supported WoL modes
	Opts      uint32 // Bitmask of WAKE_* flags for enabled WoL modes
}

// Timestamping options
// see: https://www.kernel.org/doc/Documentation/networking/timestamping.txt
const (
	SOF_TIMESTAMPING_TX_HARDWARE  = (1 << 0)  /* Request tx timestamps generated by the network adapter. */
	SOF_TIMESTAMPING_TX_SOFTWARE  = (1 << 1)  /* Request tx timestamps when data leaves the kernel. */
	SOF_TIMESTAMPING_RX_HARDWARE  = (1 << 2)  /* Request rx timestamps generated by the network adapter. */
	SOF_TIMESTAMPING_RX_SOFTWARE  = (1 << 3)  /* Request rx timestamps when data enters the kernel. */
	SOF_TIMESTAMPING_SOFTWARE     = (1 << 4)  /* Report any software timestamps when available. */
	SOF_TIMESTAMPING_SYS_HARDWARE = (1 << 5)  /* This option is deprecated and ignored. */
	SOF_TIMESTAMPING_RAW_HARDWARE = (1 << 6)  /* Report hardware timestamps. */
	SOF_TIMESTAMPING_OPT_ID       = (1 << 7)  /* Generate a unique identifier along with each packet. */
	SOF_TIMESTAMPING_TX_SCHED     = (1 << 8)  /* Request tx timestamps prior to entering the packet scheduler. */
	SOF_TIMESTAMPING_TX_ACK       = (1 << 9)  /* Request tx timestamps when all data in the send buffer has been acknowledged. */
	SOF_TIMESTAMPING_OPT_CMSG     = (1 << 10) /* Support recv() cmsg for all timestamped packets. */
	SOF_TIMESTAMPING_OPT_TSONLY   = (1 << 11) /* Applies to transmit timestamps only. */
	SOF_TIMESTAMPING_OPT_STATS    = (1 << 12) /* Optional stats that are obtained along with the transmit timestamps. */
	SOF_TIMESTAMPING_OPT_PKTINFO  = (1 << 13) /* Enable the SCM_TIMESTAMPING_PKTINFO control message for incoming packets with hardware timestamps. */
	SOF_TIMESTAMPING_OPT_TX_SWHW  = (1 << 14) /* Request both hardware and software timestamps for outgoing packets when SOF_TIMESTAMPING_TX_HARDWARE and SOF_TIMESTAMPING_TX_SOFTWARE are enabled at the same time. */
	SOF_TIMESTAMPING_BIND_PHC     = (1 << 15) /* Bind the socket to a specific PTP Hardware Clock. */
)

const (
	/*
	 * No outgoing packet will need hardware time stamping;
	 * should a packet arrive which asks for it, no hardware
	 * time stamping will be done.
	 */
	HWTSTAMP_TX_OFF = iota

	/*
	 * Enables hardware time stamping for outgoing packets;
	 * the sender of the packet decides which are to be
	 * time stamped by setting %SOF_TIMESTAMPING_TX_SOFTWARE
	 * before sending the packet.
	 */
	HWTSTAMP_TX_ON

	/*
	 * Enables time stamping for outgoing packets just as
	 * HWTSTAMP_TX_ON does, but also enables time stamp insertion
	 * directly into Sync packets. In this case, transmitted Sync
	 * packets will not received a time stamp via the socket error
	 * queue.
	 */
	HWTSTAMP_TX_ONESTEP_SYNC

	/*
	 * Same as HWTSTAMP_TX_ONESTEP_SYNC, but also enables time
	 * stamp insertion directly into PDelay_Resp packets. In this
	 * case, neither transmitted Sync nor PDelay_Resp packets will
	 * receive a time stamp via the socket error queue.
	 */
	HWTSTAMP_TX_ONESTEP_P2P
)

const (
	HWTSTAMP_FILTER_NONE                = iota /* time stamp no incoming packet at all */
	HWTSTAMP_FILTER_ALL                        /* time stamp any incoming packet */
	HWTSTAMP_FILTER_SOME                       /* return value: time stamp all packets requested plus some others */
	HWTSTAMP_FILTER_PTP_V1_L4_EVENT            /* PTP v1, UDP, any kind of event packet */
	HWTSTAMP_FILTER_PTP_V1_L4_SYNC             /* PTP v1, UDP, Sync packet */
	HWTSTAMP_FILTER_PTP_V1_L4_DELAY_REQ        /* PTP v1, UDP, Delay_req packet */
	HWTSTAMP_FILTER_PTP_V2_L4_EVENT            /* PTP v2, UDP, any kind of event packet */
	HWTSTAMP_FILTER_PTP_V2_L4_SYNC             /* PTP v2, UDP, Sync packet */
	HWTSTAMP_FILTER_PTP_V2_L4_DELAY_REQ        /* PTP v2, UDP, Delay_req packet */
	HWTSTAMP_FILTER_PTP_V2_L2_EVENT            /* 802.AS1, Ethernet, any kind of event packet */
	HWTSTAMP_FILTER_PTP_V2_L2_SYNC             /* 802.AS1, Ethernet, Sync packet */
	HWTSTAMP_FILTER_PTP_V2_L2_DELAY_REQ        /* 802.AS1, Ethernet, Delay_req packet */
	HWTSTAMP_FILTER_PTP_V2_EVENT               /* PTP v2/802.AS1, any layer, any kind of event packet */
	HWTSTAMP_FILTER_PTP_V2_SYNC                /* PTP v2/802.AS1, any layer, Sync packet */
	HWTSTAMP_FILTER_PTP_V2_DELAY_REQ           /* PTP v2/802.AS1, any layer, Delay_req packet */
	HWTSTAMP_FILTER_NTP_ALL                    /* NTP, UDP, all versions and packet modes */
)

// TimestampingInformation contains PTP timetstapming information
type TimestampingInformation struct {
	Cmd            uint32
	SoTimestamping uint32 /* SOF_TIMESTAMPING_* bitmask */
	PhcIndex       int32
	TxTypes        uint32 /* HWTSTAMP_TX_* */
	txReserved     [3]uint32
	RxFilters      uint32 /* HWTSTAMP_FILTER_ */
	rxReserved     [3]uint32
}

type EthtoolGStrings struct {
	cmd        uint32
	string_set uint32
	len        uint32
	data       [MAX_GSTRINGS * ETH_GSTRING_LEN]byte
}

type EthtoolStats struct {
	cmd     uint32
	n_stats uint32
	data    [MAX_GSTRINGS]uint64
}

// Ring is a ring config for an interface
type Ring struct {
	Cmd               uint32
	RxMaxPending      uint32
	RxMiniMaxPending  uint32
	RxJumboMaxPending uint32
	TxMaxPending      uint32
	RxPending         uint32
	RxMiniPending     uint32
	RxJumboPending    uint32
	TxPending         uint32
}

// Pause is a pause config for an interface
type Pause struct {
	Cmd     uint32
	Autoneg uint32
	RxPause uint32
	TxPause uint32
}

// Ethtool is a struct that contains the file descriptor for the ethtool
type Ethtool struct {
	fd int
}

// max values for my setup dont know how to make this dynamic
const MAX_INDIR_SIZE = 256
const MAX_CORES = 32

type Indir struct {
	Cmd       uint32
	Size      uint32
	RingIndex [MAX_INDIR_SIZE]uint32 // statically definded otherwise crash

}

type SetIndir struct {
	Equal  uint8    // used to set number of cores
	Weight []uint32 // used to select cores
}

// EthtoolCmd is the Go version of the Linux kerne ethtool_cmd struct
// see ethtool.c
type EthtoolCmd struct {
	Cmd            uint32
	Supported      uint32
	Advertising    uint32
	Speed          uint16
	Duplex         uint8
	Port           uint8
	Phy_address    uint8
	Transceiver    uint8
	Autoneg        uint8
	Mdio_support   uint8
	Maxtxpkt       uint32
	Maxrxpkt       uint32
	Speed_hi       uint16
	Eth_tp_mdix    uint8
	Reserved2      uint8
	Lp_advertising uint32
	Reserved       [2]uint32
}
