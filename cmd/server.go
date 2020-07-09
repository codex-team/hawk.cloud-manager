package cmd

import (
	"log"

	"github.com/codex-team/hawk.cloud-manager/pkg/server"
	"github.com/codex-team/hawk.cloud-manager/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{Use: "server",
		Short: "Starts cloud-manager server",
		Run: func(cmd *cobra.Command, args []string) {
			storage := storage.NewYamlStorage(storageFile)

			if err := storage.Load(); err != nil {
				log.Fatal(err)
			}

			s := server.NewServer(port, storage.Get())
			s.Run()
		},
	}
	storageFile string
	port        string
)

func init() {
	serverCmd.Flags().StringVarP(&storageFile, "storage", "s", "config.yaml", "Storage config YAML file")
	serverCmd.Flags().StringVarP(&port, "port", "p", "50051", "Listening port")

	rootCmd.AddCommand(serverCmd)
}
