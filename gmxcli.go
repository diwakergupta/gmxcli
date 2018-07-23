package main

import (
	"fmt"
	"log"

	"github.com/diwakergupta/gmxcli/cmd"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gmxcli version:", version)
	},
}

func init() {
	// Log date, time and file information by default
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cmd.RootCmd.AddCommand(versionCmd)
}

func main() {
	cmd.Execute()
}
