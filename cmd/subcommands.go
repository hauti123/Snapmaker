package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hauti123/Snapmaker/snapmaker"
)

type CliSubcommand interface {
	Init(printer snapmaker.Snapmaker) error
	Run() error
	Name() string
	Flags() *flag.FlagSet
}

var cliSubcommands = []CliSubcommand{
	NewUploadCommand(),
}

type globalFlags struct {
	printerIp       string
	apiToken        string
	discoverTimeout int
}

func runSubcommands(args []string) error {
	if len(args) < 2 {
		return errors.New("You must pass a sub-command")
	}

	globalFlags := setupGlobalFlags()

	subCommandArg := os.Args[1]
	for _, subCommand := range cliSubcommands {
		if subCommand.Name() == subCommandArg {
			subCommand.Flags().Parse(args[2:])

			printer, err := getPrinter(&globalFlags.printerIp, &globalFlags.discoverTimeout, &globalFlags.apiToken)
			if err != nil {
				return err
			}

			fmt.Printf("printer=%v\n", printer)
			err = subCommand.Init(printer)
			if err != nil {
				return err
			}
			return subCommand.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subCommandArg)
}

func setupGlobalFlags() *globalFlags {
	var globalFlags globalFlags
	for _, cliSubcommand := range cliSubcommands {
		cliSubcommand.Flags().StringVar(&globalFlags.printerIp, "printer-ip", "", "IP address of Snapmaker, automatic discovery is used if omitted.")
		cliSubcommand.Flags().StringVar(&globalFlags.apiToken, "api-token", "", "API token")
		cliSubcommand.Flags().IntVar(&globalFlags.discoverTimeout, "discover-timeout", 5, "API token")
	}
	return &globalFlags
}

func getPrinter(printerIpFlag *string, discoverTimeoutFlag *int, apiTokenFlag *string) (snapmaker.Snapmaker, error) {

	if apiTokenFlag == nil || len(*apiTokenFlag) == 0 {
		return snapmaker.Snapmaker{}, errors.New("No API token provided.")
	}

	if printerIpFlag == nil || len(*printerIpFlag) == 0 {
		if discoverTimeoutFlag == nil {
			return snapmaker.Snapmaker{}, errors.New("No discovery timeout given (should have default value).")
		}
		return snapmaker.DiscoverSnapmaker(time.Duration(*discoverTimeoutFlag)*time.Second, *apiTokenFlag)
	}

	return snapmaker.NewSnapmaker(*printerIpFlag, *apiTokenFlag), nil
}
