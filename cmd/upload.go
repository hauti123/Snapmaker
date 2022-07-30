/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload .gcode file",
	Long:  `Upload a .gcode file to the 3D printer`,
	Args:  cobra.ExactArgs(1),
	Run:   upload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func upload(cmd *cobra.Command, args []string) {
	printer, err := connectAndStoreConfig(viper.GetString(cfgPrinterIp), viper.GetInt(cfgDiscoverTimeout), viper.GetString(cfgApiToken))
	cobra.CheckErr(err)

	err = printer.SendGcodeFile(args[0])
	cobra.CheckErr(err)
}
