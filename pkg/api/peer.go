package api

import (
	"encoding/base64"
	"fmt"
	"net"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const KeyLen = wg.KeyLen

type Key [KeyLen]byte

type Peer struct {
	PublicKey  Key
	Endpoint   *string
	AllowedIPs []net.IPNet
}

type Conf struct {
	ListenPort int
	Peers      []Peer
}

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
