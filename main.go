package main

import (
	"fmt"
	"time"

	"github.com/hauti123/Snapmaker/snapmaker"
)

func statusLoop(printer *snapmaker.Snapmaker) {
	for {
		time.Sleep(1 * time.Second)
		_, _ = printer.GetStatus()
	}
}

const discoverTimeout = 5 * time.Second

func main() {
	printer, err := snapmaker.DiscoverSnapmaker(discoverTimeout)
	if err != nil {
		panic(err)
	}
	err = printer.Connect()
	if err != nil {
		panic(err)
	}

	// initial status request and status loop is needed, otherweise upload will fail with "401 - Unauthorzied, Machine not yet connected"
	status, err := printer.GetStatus()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Initial printer status:\n%s\n", status)

	go statusLoop(&printer)

	err = printer.SendGcodeFile("test.gcode")
	if err != nil {
		panic(err)
	}
}
