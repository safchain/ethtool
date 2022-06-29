package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/jpillora/opts"
	"golang.org/x/sys/unix"

	"github.com/safchain/ethtool"
)

type config struct {
	Interface        string `opts:"mode=arg,help=the name of the network device on which ethtool should operate."`
	All              bool   `opts:"help=show all information."`
	ShowFeatures     bool   `opts:"short=k,help=Queries the specified network device for the state of protocol offload and other features."`
	ShowPermAddr     bool   `opts:"short=P,help=Queries the specified network device for permanent hardware address."`
	Statistics       bool   `opts:"short=S,help=Queries the specified network device for standard or NIC- and driver-specific statistics."`
	Driver           bool   `opts:"short=i,help=Queries the specified network device for associated driver information."`
	DumpModuleEeprom bool   `opts:"short=m,help=Retrieves and if possible decodes the EEPROM from plugin modules."`
	ShowRxfhIndir    bool   `opts:"short=x,help=Retrieves the receive flow hash indirection table and/or RSS hash key."`
	SetRxfhIndir     bool   `opts:"short=X,help=Configures the receive flow hash indirection table and/or RSS hash key."`
	HashKey          []byte `opts:"group=set-rxfh-indir,help=Sets RSS hash key of the specified network device."`
	HashFunc         string `opts:"group=set-rxfh-indir,help=Sets RSS hash function of the specified network device."`
	Start            int    `opts:"group=set-rxfh-indir,help=For the equal and weight options sets the starting receive queue for spreading flows to N."`
	Equal            int    `opts:"group=set-rxfh-indir,help=Sets the receive flow hash indirection table to spread flows evenly between the first N receive queues."`
	Weight           []int  `opts:"group=set-rxfh-indir,help=Sets the receive flow hash indirection table to spread flows between receive queues according to the given weights."`
	Default          bool   `opts:"group=set-rxfh-indir,help=Sets the receive flow hash indirection table to its default value."`
	Context          int    `opts:"group=set-rxfh-indir,help=Specifies an RSS context to act on."`
	Delete           bool   `opts:"group=set-rxfh-indir,help=Delete the specified RSS context."`
}

func main() {
	c := new(config)
	opts.Parse(c)

	e, err := ethtool.NewEthtool()
	if err != nil {
		panic(err.Error())
	}
	defer e.Close()

	if c.ShowFeatures || c.All {
		if err = showFeatures(e, c); err != nil {
			panic(err)
		}
	}

	if c.ShowPermAddr || c.All {
		if err = showPermAddr(e, c); err != nil {
			panic(err)
		}
	}

	if c.Statistics || c.All {
		if err = showStats(e, c); err != nil {
			panic(err)
		}
	}

	if c.Driver || c.All {
		if err = showDriver(e, c); err != nil {
			panic(err)
		}
	}

	if c.DumpModuleEeprom || c.All {
		if err = dumpModuleEeprom(e, c); err != nil {
			panic(err)
		}
	}

	if c.ShowRxfhIndir || c.All {
		if err = showRxfhIndir(e, c); err != nil {
			panic(err)
		}
	}

	if !(c.ShowFeatures || c.ShowPermAddr || c.Statistics || c.Driver || c.DumpModuleEeprom || c.ShowRxfhIndir) || c.All {
		if err = showSettings(e, c); err != nil {
			panic(err)
		}
	}
}

func showFeatures(e *ethtool.Ethtool, c *config) error {
	features, err := e.Features(c.Interface)
	if err != nil {
		return err
	}
	fmt.Println("Features for ", c.Interface)

	var keys []string

	for name := range features {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		on := features[name]
		status := "off"
		if on {
			status = "on"
		}
		fmt.Printf("%s: %s\n", name, status)
	}

	return nil
}

func showPermAddr(e *ethtool.Ethtool, c *config) error {
	permAddr, err := e.PermAddr(c.Interface)
	if err != nil {
		return err
	}

	fmt.Println("Permanent address:", permAddr)

	return nil
}

func showStats(e *ethtool.Ethtool, c *config) error {
	stats, err := e.Stats(c.Interface)
	if err != nil {
		return err
	}
	fmt.Println("NIC statistics")

	var keys []string

	for name := range stats {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	for _, name := range keys {
		fmt.Printf("\t%s: %+v\n", name, stats[name])
	}

	return nil
}

func showDriver(e *ethtool.Ethtool, c *config) error {
	drvr, err := e.DriverName(c.Interface)
	if err != nil {
		return err
	}
	fmt.Println("driver:", drvr)

	drvInfo, err := e.DriverInfo(c.Interface)
	if err != nil {
		return err
	}
	fmt.Println("version:", drvInfo.Version)
	fmt.Println("firmware-version:", drvInfo.FwVersion)
	fmt.Println("expansion-rom-version:", drvInfo.EromVersion)
	fmt.Println("bus-info:", drvInfo.BusInfo)
	fmt.Println("supports-statistics:", drvInfo.NStats != 0)
	fmt.Println("supports-test:", drvInfo.TestInfoLen != 0)
	fmt.Println("supports-eeprom-access:", drvInfo.EedumpLen != 0)
	fmt.Println("supports-register-dump:", drvInfo.RegdumpLen != 0)
	fmt.Println("supports-priv-flags:", drvInfo.NPrivFlags != 0)

	return nil
}

func dumpModuleEeprom(e *ethtool.Ethtool, c *config) error {
	eeprom, err := e.ModuleEepromHex(c.Interface)
	if err != nil {
		if errors.Is(err, syscall.ENOTSUP) || errors.Is(err, syscall.EPERM) {
			fmt.Fprintln(os.Stderr, "Cannot get module EEPROM information:", err)
		} else {
			return err
		}
	} else {
		fmt.Printf("module eeprom: %+v\n", eeprom)
	}

	return nil
}

func showRxfhIndir(e *ethtool.Ethtool, c *config) error {
	rssHash, err := e.GetFlowHash(c.Interface)
	if err != nil {
		if errors.Is(err, syscall.ENOTSUP) || errors.Is(err, syscall.EPERM) {
			fmt.Fprintf(os.Stderr, "RX flow hash: %s\n", err)
		} else {
			return err
		}
	} else {
		fmt.Printf("RX flow hash indirection table for %s with %d RX ring(s):\n", c.Interface, rssHash.RingCount)
		if len(rssHash.Table) == 0 {
			fmt.Println("Operation not supported")
		} else {
			fmt.Println(rssHash.Table)
		}

		fmt.Println("RSS hash key:")
		if len(rssHash.Key) == 0 {
			fmt.Println("Operation not supported")
		} else {
			fmt.Println(hex.EncodeToString(rssHash.Key))
		}

		fmt.Println("RSS hash function:")
		for name, b := range rssHash.Funcs {
			status := "off"
			if b {
				status = "on"
			}
			fmt.Printf("    %s: %s\n", name, status)
		}
	}

	return err
}

func showSettings(e *ethtool.Ethtool, c *config) error {
	m, err := e.CmdGetMapped(c.Interface)
	if err != nil {
		return err
	}

	msgLvl, err := e.MsglvlGet(c.Interface)
	if err != nil {
		return err
	}

	fmt.Println("Settings for", c.Interface)
	dumpLinkCaps("Supported", []uint32{uint32(m["Supported"])})
	dumpLinkCaps("Advertisied", []uint32{uint32(m["Advertising"])})
	fmt.Println("\tSpeed:", m["speed"], "Mb/s")
	fmt.Println("\tDuplex:", ethtool.DuplexName(uint8(m["Duplex"])))
	fmt.Println("\tPort:", ethtool.PortTypeName(uint8(m["Port"])))
	fmt.Println("\tPHYAD:", m["Phy_address"])
	fmt.Println("\tTransceiver:", ethtool.TransceiverName(uint8(m["Transceiver"])))

	status := "off"
	if m["Autoneg"] != 0 {
		status = "on"
	}
	fmt.Println("\tAuto-negotiation::", status)

	if m["Port"] == ethtool.PORT_TP {
		dumpMdix(uint8(m["Eth_tp_mdix"]), uint8(m["Eth_tp_mdix_ctrl"]))
	}

	fmt.Printf("\tCurrent message level: 0x%08x (%d)\n", msgLvl, msgLvl)
	fmt.Println("\t                      ", strings.Join(ethtool.MsgLevelNames(msgLvl), " "))

	return nil
}

func dumpLinkCaps(prefix string, v []uint32) {
	fmt.Printf("\t%s ports: [ %s ]\n", prefix, strings.Join(ethtool.LinkPortNames(v), " "))
	fmt.Printf("\t%s link modes: %s\n", prefix, strings.Join(ethtool.LinkSpeedNames(v), " "))
	fmt.Printf("\t%s pause frame use: ", prefix)
	if ethtool.LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_Pause_BIT) {
		fmt.Printf("Symmetric")
		if ethtool.LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_Asym_Pause_BIT) {
			fmt.Printf(" Receive-only")
		}
		fmt.Println()
	} else {
		if ethtool.LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_Asym_Pause_BIT) {
			fmt.Println("Transmit-only")
		} else {
			fmt.Println("No")
		}
	}
	fmt.Printf("\t%s auto-negotiation: %t\n", prefix, ethtool.LinkModeTestBit(v, unix.ETHTOOL_LINK_MODE_Autoneg_BIT))
	fmt.Printf("\t%s FEC modes: ", prefix)
	if s := ethtool.LinkECCModeNames(v); len(s) > 0 {
		fmt.Println(strings.Join(s, " "))
	} else {
		fmt.Println("Not reported")
	}
}

func dumpMdix(mdix, mdix_ctrl uint8) {
	fmt.Printf("\tMDI-X: ")
	switch mdix_ctrl {
	case ethtool.ETH_TP_MDI:
		fmt.Println("off (forced)")
	case ethtool.ETH_TP_MDI_X:
		fmt.Println("on (forced)")
	default:
		switch mdix {
		case ethtool.ETH_TP_MDI:
			fmt.Printf("off")
		case ethtool.ETH_TP_MDI_X:
			fmt.Printf("on")
		default:
			fmt.Printf("unknown")
		}

		if mdix_ctrl == ethtool.ETH_TP_MDI_AUTO {
			fmt.Printf(" (auto)")
		}

		fmt.Println()
	}
}
