package main

import (
	"fmt"
	"os"
)

func main() {
	err := runSubcommands(os.Args)

	if err != nil {
		fmt.Printf("Fail: %v", err)
	}
}
