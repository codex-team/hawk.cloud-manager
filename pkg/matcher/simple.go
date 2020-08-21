package matcher

import (
	"fmt"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
)

type Simple struct {
	config     config.PeerConfig
	hostsKeys  map[api.Key]*config.Host
	hostsNames map[string]*config.Host
}

func NewSimpleMatcher(cfg config.PeerConfig) (*Simple, error) {
	var m = &Simple{config: cfg}
	m.hostsKeys = make(map[api.Key]*config.Host)
	m.hostsNames = make(map[string]*config.Host)
	for _, host := range cfg.Hosts {
		key, err := api.NewKey(host.PublicKey)
		if err != nil {
			return nil, err
		}
		m.hostsKeys[key] = &host
		m.hostsNames[host.Name] = &host
	}
	return m, nil
}

// Peers to connect to
func (m *Simple) Peers(key api.Key) ([]api.Peer, error) { // TODO: it's so bad...
	var peers []api.Peer
	sourcePeer, ok := m.hostsKeys[key]
	if !ok {
		return []api.Peer{}, nil
	}
	for _, group := range m.config.Groups {
		for _, hostName := range group.Hosts {
			host, ok := m.hostsNames[hostName]
			if ok && sourcePeer == host {
				for _, futurePeerHostName := range group.Hosts {
					futurePeer, ok := m.hostsNames[futurePeerHostName]
					if !ok {
						return nil, fmt.Errorf("can't find peer from group") // TODO: move to storage validation
					}
					apiPeer, err := futurePeer.ToAPIPeer()
					if err != nil {
						return nil, err
					}
					peers = append(peers, *apiPeer)
				}
			}
		}
	}
	return peers, nil
}
