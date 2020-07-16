package cmd

import (
	"context"
	"github.com/codex-team/hawk.cloud-manager/pkg/server"
	"github.com/codex-team/hawk.cloud-manager/pkg/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serverCmd = &cobra.Command{Use: "server",
		Short: "Starts cloud-manager server",
		Run: func(cmd *cobra.Command, args []string) {
			logger, _ := zap.NewDevelopment()
			defer logger.Sync()
			sugared := logger.Sugar()

			storage := storage.NewYamlStorage(storageFile)

			if err := storage.Load(); err != nil {
				sugared.Fatal(err)
			}

			manager, err := server.New(addr, storage.Get(), *logger)
			if err != nil {
				sugared.Fatal(err)
			}

			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			srv := &http.Server{
				Addr: addr,
				Handler: manager.Handler,
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					sugared.Fatalw("Server listen error", err)
				}
			}()
			sugared.Infof("Server started at %s", addr)

			<-done
			sugared.Info("Server stopped")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer func() {
				//h.Shutdown()
				cancel()
			}()

			if err := srv.Shutdown(ctx); err != nil {
				sugared.Fatalf("Server shutdown failed:%+v", err)
			}
			sugared.Info("Server exited")

		},
	}
	storageFile string
	addr        string
)

func init() {
	serverCmd.Flags().StringVarP(&storageFile, "storage", "s", "config.yaml", "Storage config YAML file")
	serverCmd.Flags().StringVarP(&addr, "addr", "l", "0.0.0.0:50051", "Listening addr")

	rootCmd.AddCommand(serverCmd)
}
