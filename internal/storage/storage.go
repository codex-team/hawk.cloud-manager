package storage

import "github.com/codex-team/hawk.cloud-manager/pkg/config"

// Storage stores peer configuration
type Storage interface {
	Load() error
	Get() *config.PeerConfig
}
