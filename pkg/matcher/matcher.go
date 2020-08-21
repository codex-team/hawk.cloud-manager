package matcher

import "github.com/codex-team/hawk.cloud-manager/pkg/api"

type Matcher interface {
	Peers(key api.Key) ([]api.Peer, error)
}
