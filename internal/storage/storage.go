package storage

import "github.com/codex-team/hawk.cloud-manager/pkg/config"

// Storage stores peer configuration
type Storage interface {
	// Load reads Peer Config from file and saves it to `config` field
	Load() error
	// Save writes Peer Config from `config` field to file
	Save() error
	// Get returns stored Peer Config
	Get() *config.PeerConfig
	// Set updates stored Peer Config
	Set(config.PeerConfig)
}
