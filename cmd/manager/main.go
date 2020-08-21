package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/codex-team/hawk.cloud-manager/internal/server"
	"github.com/codex-team/hawk.cloud-manager/internal/storage/yaml"
	"go.uber.org/zap"
)

var (
	storageFile string
	addr        string
)

func init() {
	flag.StringVar(&storageFile, "storage", "examples/config.yaml", "Storage config YAML file")
	flag.StringVar(&addr, "addr", "0.0.0.0:50051", "Listening addr")
}

func main() {
	flag.Parse()
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugared := logger.Sugar()

	storage := yaml.NewYamlStorage(storageFile)
	if err := storage.Load(); err != nil {
		sugared.Fatal(err)
	}

	manager, err := server.New(addr, storage.Get(), logger)
	if err != nil {
		sugared.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := manager.Stop(); err != nil {
			sugared.Error("server stopped with error", zap.Error(err))
			return
		}
	}()

	go func() {
		sugared.Infof("server started at %s", addr)
		errs <- manager.Run()
	}()

	select {
	case <-done:
		signal.Stop(done)
		return
	case err = <-errs:
		if err != nil {
			sugared.Error("server exited with error", zap.Error(err))
		}
		return
	}
}
