package config

import (
	"net"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
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
func (p *PeerConfig) ToAPIConf() (*api.Conf, error) {
	cf := api.Conf{
		Peers: []api.Peer{},
	}

	for _, h := range p.Hosts {
		err := api.CheckKeyFormat(h.PublicKey)
		if err != nil {
			return nil, err
		}
		ips := make([]net.IPNet, len(h.AllowedIPs))
		for i, cidr := range h.AllowedIPs {
			_, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, err
			}
			ips[i] = *ipnet
		}

		cf.Peers = append(cf.Peers, api.Peer{
			PublicKey:  h.PublicKey,
			Endpoint:   h.Endpoint,
			AllowedIPs: ips,
		})
	}

	return &cf, nil
}
