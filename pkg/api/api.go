package api

import (
	"encoding/base64"
	"fmt"
	"time"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const KeyLen = wg.KeyLen

// Creds is request data required to get WireGuard configuration
type Creds struct {
	// WireGuard public key
	PublicKey string `json:"public_key"`
	// Ed25519 signature
	Signature string `json:"signature"`
}

type Peer struct {
	// WireGuard peer public key (a base64-encoded string that is 32 bytes long)
	PublicKey string `json:"public_key" yaml:"public_key"`
	// WireGuard peer endpoint
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	// WireGuard peer keep alive interval
	PersistentKeepAliveInterval time.Duration `json:"keep_alive_interval,omitempty" yaml:"keep_alive_interval,omitempty"`
	// WireGuard peer allowed IPs
	AllowedIPs []string `json:"allowed_ips" yaml:"allowed_ips"`
}

// Conf is WireGuard configuration
type Conf struct {
	// WireGuard Listen port
	ListenPort int `json:"listen_port" yaml:"listen_port"`
	// Peers list
	Peers []Peer `json:"peers" yaml:"peers"`
}

// Equals compares this Conf to other
func (c *Conf) Equals(other *Conf) bool {
	if c.ListenPort != other.ListenPort {
		return false
	}
	if len(c.Peers) != len(other.Peers) {
		return false
	}
	for i, p := range c.Peers {
		switch {
		case p.PublicKey != other.Peers[i].PublicKey:
			return false
		case p.Endpoint != other.Peers[i].Endpoint:
			return false
		case p.PersistentKeepAliveInterval != other.Peers[i].PersistentKeepAliveInterval:
			return false
		case len(p.AllowedIPs) != len(other.Peers[i].AllowedIPs):
			return false
		default:
			for i, ip := range other.Peers[i].AllowedIPs {
				if ip != p.AllowedIPs[i] {
					return false
				}
			}
		}
	}

	return true
}

// CheckKeyFormat checks if public key is valid
func CheckKeyFormat(s string) error {
	key, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("failed to parse base64-encoded key: %v", err)
	}
	if len(key) != KeyLen {
		return fmt.Errorf("incorrect key size: %d", len(key))
	}

	return nil
}
