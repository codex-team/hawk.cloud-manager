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
	PublicKey                   Key           `json:"public_key"`
	Endpoint                    string        `json:"endpoint"`
	PersistentKeepAliveInterval time.Duration `json:"keep_alive_interval"`
	AllowedIPs                  []net.IPNet   `json:"allowed_ips"`
}

// Conf is WireGuard configuration
type Conf struct {
	ListenPort int    `json:"listen_port"`
	Peers      []Peer `json:"peers"`
}

// Creds is request data required to get WireGuard configuration
type Creds struct {
	PublicKey string `json:"public_key"`
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

func (k Key) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}
