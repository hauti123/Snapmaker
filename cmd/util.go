package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/hauti123/Snapmaker/snapmaker"
	"github.com/spf13/viper"
)

func findPrinter(printerIp string, discoverTimeout int, apiToken string) (*snapmaker.Snapmaker, error) {
	var printer *snapmaker.Snapmaker
	var err error

	if len(printerIp) == 0 {
		printer, err = discoverPrinter(discoverTimeout, apiToken)
	} else {
		printer, err = createPrinter(printerIp, apiToken)
	}

	if err != nil {
		return nil, err
	}

	return printer, nil
}

func connectAndStoreConfig(printerIp string, discoverTimeout int, apiToken string) (*snapmaker.Snapmaker, error) {
	printer, err := findPrinter(printerIp, discoverTimeout, apiToken)

	if err != nil {
		return nil, err
	}

	err = printer.Connect()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Please confirm connection on Snapmaker controller")
	err = printer.WaitForConnection(60 * time.Second)
	if err != nil {
		return nil, err
	}

	storePrinterConfig(printer)

	return printer, nil
}

func discoverPrinter(discoverTimeout int, apiToken string) (*snapmaker.Snapmaker, error) {
	if discoverTimeout == 0 {
		return nil, errors.New("No discovery timeout given (should have default value).")
	}
	return snapmaker.DiscoverSnapmaker(time.Duration(discoverTimeout)*time.Second, apiToken)
}

func createPrinter(printerIp string, apiToken string) (*snapmaker.Snapmaker, error) {
	return snapmaker.NewSnapmaker(printerIp, apiToken), nil
}

func storePrinterConfig(printer *snapmaker.Snapmaker) {

	fmt.Printf("storing configuration to %s\n", viper.GetViper().ConfigFileUsed())

	viper.Set(cfgPrinterIp, printer.GetIpAdress())
	viper.Set(cfgApiToken, printer.GetApiToken())
	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("failed writing config: %v\n", err)
	}
}
