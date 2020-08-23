package config

import (
	"net"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
)

type Host struct {
	Name       string   `yaml:"name" json:"name"`
	PublicKey  string   `yaml:"public_key" json:"public_key"`
	Endpoint   string   `yaml:"endpoint" json:"endpoint,omitempty"`
	AllowedIPs []string `yaml:"allowed_ips" json:"allowed_ips"`
}

type Group struct {
	Name  string   `yaml:"name" json:"name"`
	Hosts []string `yaml:"hosts" json:"hosts"`
}

type PeerConfig struct {
	Hosts  []Host  `yaml:"hosts" json:"hosts"`
	Groups []Group `yaml:"groups" json:"groups"`
}

func (p *PeerConfig) ToAPIConf() (*api.Conf, error) {
	cf := api.Conf{
		Peers: []api.Peer{},
	}

	for _, h := range p.Hosts {
		key, err := api.NewKey(h.PublicKey)
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
			PublicKey:  key,
			Endpoint:   h.Endpoint,
			AllowedIPs: ips,
		})
	}

	return &cf, nil
}
