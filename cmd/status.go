/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print current printer status",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: printStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printStatus(cmd *cobra.Command, args []string) {
	printer, err := connectAndStoreConfig(viper.GetString(cfgPrinterIp), viper.GetInt(cfgDiscoverTimeout), viper.GetString(cfgApiToken))
	cobra.CheckErr(err)

	status, err := printer.GetStatus(time.Duration(viper.GetInt(cfgDiscoverTimeout)) * time.Second)

	yamlStatus, err := yaml.Marshal(status)
	cobra.CheckErr(err)

	fmt.Printf("Printer status:\n%s\n", yamlStatus)
}
