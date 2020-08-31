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
	// Storage config YAML file
	storageFile string
	// Address to listen
	addr string
)

// init initializes command line flags
func init() {
	flag.StringVar(&storageFile, "config", "examples/config.yaml", "Storage config YAML file")
	flag.StringVar(&addr, "addr", "0.0.0.0:50051", "Listening addr")
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Create Storage and load Config
	storage := yaml.NewYamlStorage(storageFile)
	if err := storage.Load(); err != nil {
		log.Fatal(err)
	}

	// Create server
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

	// Run server
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
