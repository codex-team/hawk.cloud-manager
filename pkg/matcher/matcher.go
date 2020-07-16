package matcher

import "github.com/codex-team/hawk.cloud-manager/api"

type Matcher interface {
	Peers(key api.Key) ([]api.Peer, error)
}
