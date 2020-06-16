package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{Use: "server",
	Short: "Starts cloud-manager server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Server")
	},
}
