/*
Copyright © 2022 Stefan Hauth snap@hauth.at

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgPrinterIp       = "ip"
	cfgDiscoverTimeout = "discoverTimeout"
	cfgApiToken        = "apiToken"
)

type rootCmdLineArguments struct {
	configFile           string
	printerIp            string
	apiToken             string
	discoverTimeoutS     int
	confirmationTimeoutS int
}

var (
	rootCmdArgs rootCmdLineArguments

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "snap",
		Short: "Snapmaker CLI tool",
		Long: `Snap is a CLI tool that allows communication with
	Snapmaker 3D printers. It currently supports .gcode file upload.`,

		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.configFile, "config", "c", "", "config file (default is $HOME/.snap.yaml)")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.printerIp, "ip", "", "", "IP address of printer. Auto discover will be used if left empty.")
	rootCmd.PersistentFlags().StringVarP(&rootCmdArgs.apiToken, "token", "t", "", "API token for communication with the 3D printer (required).")
	rootCmd.PersistentFlags().IntVarP(&rootCmdArgs.discoverTimeoutS, "discoverTimeout", "", 5, "Timeout for auto discover of 3D printer in seconds.")
	rootCmd.PersistentFlags().IntVarP(&rootCmdArgs.confirmationTimeoutS, "confirmationTimeout", "", 60, "Timeout until a new connection has to be confirmed on the Snapmaker controller (in seconds).")

	viper.BindPFlag(cfgPrinterIp, rootCmd.PersistentFlags().Lookup("ip"))
	viper.BindPFlag(cfgApiToken, rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag(cfgDiscoverTimeout, rootCmd.PersistentFlags().Lookup("discoverTimeout"))

	// Cobra also supportslocal flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootCmdArgs.configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(rootCmdArgs.configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".snap" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".snap")
		viper.SetConfigFile(fmt.Sprintf("%s/.snap.yaml", home))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.WriteConfig(); err != nil {
		fmt.Printf("failed writing config: %v\n", err)
	}
}
