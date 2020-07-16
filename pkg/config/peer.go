package config

import (
	"github.com/codex-team/hawk.cloud-manager/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/utils"
	"net"
)

type Host struct {
	Name      string `yaml:"name" json:"name"`
	PublicKey string `yaml:"public_key" json:"public_key"`
	Endpoint string `yaml:"endpoint" json:"endpoint,omitempty"`
	AllowedIPs []string `yaml:"allowed_ips" json:"allowed_ips"`
}

type Group struct {
	Name string `yaml:"name" json:"name"`
	Hosts []string `yaml:"hosts" json:"hosts"`
}

type PeerConfig struct {
	Hosts  []Host `yaml:"hosts" json:"hosts"`
	Groups []Group `yaml:"groups" json:"groups"`
}

func (h *Host)ToAPIPeer() (*api.Peer, error) {
	key, err := utils.ParseKey(h.PublicKey)
	if err != nil {
		return nil, err
	}
	var peer = api.Peer{
		PublicKey:  key,
		Endpoint:   &h.Endpoint,
	}
	peer.AllowedIPs = make([]net.IPNet, len(h.AllowedIPs))
	for i, cidr := range h.AllowedIPs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		peer.AllowedIPs[i] = *ipnet
	}
	return &peer, nil
}