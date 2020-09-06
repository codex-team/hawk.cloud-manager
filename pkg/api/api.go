package api

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const KeyLen = wg.KeyLen

// Key is base64-encoded byte array that is 32 bytes long
type Key [KeyLen]byte

type Peer struct {
	// WireGuard peer public key
	PublicKey Key `json:"public_key"`
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

// NewKey parses a Key from a string
func NewKey(s string) (Key, error) {
	key := Key{}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return key, fmt.Errorf("failed to parse base64-encoded key: %v", err)
	}
	if len(b) != KeyLen {
		return key, fmt.Errorf("incorrect key size: %d", len(b))
	}
	copy(key[:], b)

	return key, nil
}

// String returns string representation of Key
func (k Key) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}
