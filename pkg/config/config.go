package config

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
)

var (
	ErrNoPeers          = errors.New("no peers found in config")
	ErrUnknownPublicKey = errors.New("unknown public key")
)

// Host is information about a host included in PeerConfig
type Host struct {
	Name       string   `yaml:"name" json:"name"`
	PublicKey  string   `yaml:"public_key" json:"public_key"`
	Endpoint   string   `yaml:"endpoint" json:"endpoint,omitempty"`
	AllowedIPs []string `yaml:"allowed_ips" json:"allowed_ips"`
}

// Group is association of hosts
type Group struct {
	Name  string   `yaml:"name" json:"name"`
	Hosts []string `yaml:"hosts" json:"hosts"`
}

// PeerConfig is peer configuration stored in Storage
type PeerConfig struct {
	Hosts  []Host  `yaml:"hosts" json:"hosts"`
	Groups []Group `yaml:"groups" json:"groups"`
}

// ToAPIConf converts PeerConfig to api.Conf
func (p *PeerConfig) ToAPIConf(publicKey string) (*api.Conf, error) {
	cf := api.Conf{
		Peers: []api.Peer{},
	}

	publicKeyFound := false
	for _, h := range p.Hosts {
		err := api.CheckKeyFormat(h.PublicKey)
		if err != nil {
			return nil, err
		}
		if h.PublicKey == publicKey {
			publicKeyFound = true
			port := strings.Split(h.Endpoint, ":")
			if len(port) == 2 {
				listenPort, err := strconv.ParseInt(port[1], 10, 64)
				if err != nil {
					return nil, err
				}
				cf.ListenPort = int(listenPort)
			}
			continue
		}
		for _, cidr := range h.AllowedIPs {
			_, _, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, err
			}
		}

		cf.Peers = append(cf.Peers, api.Peer{
			PublicKey:  h.PublicKey,
			Endpoint:   h.Endpoint,
			AllowedIPs: h.AllowedIPs,
		})
	}

	if !publicKeyFound {
		return nil, ErrUnknownPublicKey
	}

	if len(cf.Peers) == 0 {
		return nil, ErrNoPeers
	}

	return &cf, nil
}
