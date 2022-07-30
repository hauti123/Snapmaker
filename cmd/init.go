/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the discover command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Establish a connection to a Snapmaker 3D printer and store configuration.",
	Run:   initPrinterConnection,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initPrinterConnection(cmd *cobra.Command, args []string) {
	_, err := connectAndStoreConfig(viper.GetString(cfgPrinterIp), viper.GetInt(cfgDiscoverTimeout), viper.GetString(cfgApiToken))
	cobra.CheckErr(err)
}
