package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codex-team/hawk.cloud-manager/internal/server"
	"github.com/codex-team/hawk.cloud-manager/internal/storage/yaml"
)

var (
	storageFile string
	addr        string
)

func init() {
	flag.StringVar(&storageFile, "config", "examples/config.yaml", "Storage config YAML file")
	flag.StringVar(&addr, "addr", "0.0.0.0:50051", "Listening addr")
}

func main() {
	flag.Parse()

	storage := yaml.NewYamlStorage(storageFile)
	if err := storage.Load(); err != nil {
		log.Fatal(err)
	}

	manager, err := server.New(addr, storage.Get())
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := manager.Stop(); err != nil {
			log.Fatal("server stopped with error: %w", err)
			return
		}
	}()

	go func() {
		log.Printf("server started at %s", addr)
		errs <- manager.Run()
	}()

	select {
	case <-done:
		signal.Stop(done)
		return
	case err = <-errs:
		if err != nil {
			log.Fatal("server exited with error: %w", err)
		}
		return
	}
}
