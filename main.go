package main

import (
	"fmt"
	"time"

	"github.com/hauti123/Snapmaker/snapmaker"
)

func statusLoop(printerIp string) {
	for {
		time.Sleep(1 * time.Second)
		_, _ = snapmaker.GetPrinterStatus(printerIp)
	}
}

const discoverTimeout = 5 * time.Second

func main() {
	printerIp, err := snapmaker.DiscoverSnapmaker(discoverTimeout)
	if err != nil {
		panic(err)
	}
	err = snapmaker.ConnectToPrinter(printerIp)
	if err != nil {
		panic(err)
	}

	// initial status request is needed, otherweise upload will fail with "401 - Unauthorzied, Machine not yet connected"
	status, err := snapmaker.GetPrinterStatus(printerIp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Initial printer status:\n%s\n", status)

	go statusLoop(printerIp)
	err = snapmaker.SendGcodeFile(printerIp, "test.gcode")
	if err != nil {
		panic(err)
	}
}
