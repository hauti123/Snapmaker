package main

import (
	"fmt"
	"time"
)

func statusLoop() {
	for {
		time.Sleep(1 * time.Second)
		_ = status()
	}
}

func main() {
	discoverSnapmaker()

	connect()
	// initial status request is needed, otherweise upload will fail with "401 - Unauthorzied, Machine not yet connected"
	fmt.Printf("initial status: %s\n\n", status())
	//	go statusLoop()
	sendFile("test.gcode")
}
