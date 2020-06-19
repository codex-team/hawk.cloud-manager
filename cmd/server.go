package cmd

import (
	"github.com/spf13/cobra"
	"github.com/codex-team/hawk.cloud-manager/pkg/server"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{Use: "server",
	Short: "Starts cloud-manager server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}
