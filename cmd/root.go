package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "gmxcli",
	SilenceUsage: true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by gmxcli.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
}
