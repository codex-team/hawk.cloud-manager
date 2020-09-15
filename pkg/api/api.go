package api

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const KeyLen = wg.KeyLen

type Peer struct {
	// WireGuard peer public key (a base64-encoded string that is 32 bytes long)
	PublicKey string `json:"public_key"`
	// WireGuard peer endpoint
	Endpoint string `json:"endpoint,omitempty"`
	// WireGuard peer keep alive interval
	PersistentKeepAliveInterval time.Duration `json:"keep_alive_interval,omitempty"`
	// WireGuard peer allowed IPs
	AllowedIPs []net.IPNet `json:"allowed_ips"`
}

// Conf is WireGuard configuration
type Conf struct {
	// WireGuard Listen port
	ListenPort int `json:"listen_port"`
	// Peers list
	Peers []Peer `json:"peers"`
}

// Creds is request data required to get WireGuard configuration
type Creds struct {
	// WireGuard public key
	PublicKey string `json:"public_key"`
	// Ed25519 signature
	Signature string `json:"signature"`
}

// ParseKey checks if public key is valid
func ParseKey(s string) error {
	key, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("failed to parse base64-encoded key: %v", err)
	}
	if len(key) != KeyLen {
		return fmt.Errorf("incorrect key size: %d", len(key))
	}

	return nil
}
