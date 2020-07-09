package cmd

import (
	"github.com/codex-team/hawk.cloud-manager/pkg/agent"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(agentCmd)
}

var agentCmd = &cobra.Command{Use: "agent",
	Short: "Starts cloud-manager agent",
	Run: func(cmd *cobra.Command, args []string) {
		agent.Run()
	},
}

