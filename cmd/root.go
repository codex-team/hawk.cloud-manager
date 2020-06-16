package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cloud-manager",
	Short: "cloud-manager synchronizes Wireguard peer config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("It's alive!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}