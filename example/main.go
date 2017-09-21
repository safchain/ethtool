package main

import (
	"flag"
	"fmt"
	"log"

	//	"github.com/safchain/ethtool"
	"github.com/mickep76/ethtool"
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

	eeprom, err := e.ModuleEepromHex(*name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("module eeprom: %+v\n", eeprom)
}
