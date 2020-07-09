package storage

import "github.com/codex-team/hawk.cloud-manager/pkg/config"

type Storage interface {
	Load() error
	Get() config.PeerConfig
}
