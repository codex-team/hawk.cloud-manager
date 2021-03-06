package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/codex-team/hawk.cloud-manager/internal/agent"
)

var (
	// CloudManager address
	managerAddress string
	// Config file path
	configFile string
	// Time interval to check config changes
	interval time.Duration
	// Path to file with public key
	pubKeyFile string
	// Path to file with private key
	privKeyFile string
)

// init initializes command line flags
func init() {
	flag.StringVar(&configFile, "config", "/etc/wireguard/wg0.conf", "file to store WireGuard config")
	flag.StringVar(&managerAddress, "manager", "http://manager:50051", "CloudManager address")
	flag.StringVar(&pubKeyFile, "pubkey", "agent_pubKey", "path to file with public key")
	flag.StringVar(&privKeyFile, "privkey", "agent_privKey", "path to file with private key")
	flag.DurationVar(&interval, "interval", time.Minute, "time interval to check config changes")
}

// formatKey returns key string without newlines
func formatKey(key []byte) string {
	return strings.Replace(strings.Replace(string(key), "\n", "", -1), " ", "", -1)
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Read public key
	pubKeyData, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Read private key
	privKeyData, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create agent
	agentInstance, err := agent.New(managerAddress, configFile, formatKey(pubKeyData), formatKey(privKeyData))
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("starting agent")

	err = agentInstance.PullNewConf()
	if err != nil {
		log.Fatal(err)
	}

	// Query new config from Manager periodically
	ticker := time.NewTicker(interval)
	stop := false
	for !stop {
		select {
		case <-done:
			signal.Stop(done)
			stop = true
		case err = <-errs:
			if err != nil {
				log.Fatal(err)
			}
			stop = true
		case <-ticker.C:
			errs <- agentInstance.PullNewConf()
		}
	}
}
