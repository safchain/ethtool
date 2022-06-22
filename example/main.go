package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/safchain/ethtool"
)

func main() {
	name := flag.String("interface", "", "Interface name")
	flag.Parse()

	if *name == "" {
		log.Fatal("interface is not specified")
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		panic(err.Error())
	}
	defer e.Close()

	features, err := e.Features(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("features: %+v\n", features)

	stats, err := e.Stats(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("stats: %+v\n", stats)

	busInfo, err := e.BusInfo(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("bus info: %+v\n", busInfo)

	drvr, err := e.DriverName(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("driver name: %+v\n", drvr)

	cmdGet, err := e.CmdGetMapped(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("cmd get: %+v\n", cmdGet)

	msgLvlGet, err := e.MsglvlGet(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("msg lvl get: %+v\n", msgLvlGet)

	drvInfo, err := e.DriverInfo(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("drvrinfo: %+v\n", drvInfo)

	permAddr, err := e.PermAddr(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("permaddr: %+v\n", permAddr)

	eeprom, err := e.ModuleEepromHex(*name)
	if err != nil {
		if errors.Is(err, syscall.ENOTSUP) || errors.Is(err, syscall.EPERM) {
			fmt.Fprintf(os.Stderr, "module eeprom: %s\n", err)
		} else {
			panic(err.Error())
		}
	} else {
		fmt.Printf("module eeprom: %+v\n", eeprom)
	}

	rssHash, err := e.GetFlowHash(*name)
	if err != nil {
		if errors.Is(err, syscall.ENOTSUP) || errors.Is(err, syscall.EPERM) {
			fmt.Fprintf(os.Stderr, "RX flow hash: %s\n", err)
		} else {
			panic(err.Error())
		}
	} else {
		fmt.Printf("RX flow hash indirection table for %s with %d RX ring(s):\n", *name, rssHash.RingCount)
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
}
