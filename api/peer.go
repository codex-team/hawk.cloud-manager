package api

import (
	"encoding/base64"
	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
)

const KeyLen = wg.KeyLen

type Key [KeyLen]byte

type Peer struct {
	PublicKey Key
	Endpoint *string
	AllowedIPs []net.IPNet
}

type Conf struct {
	ListenPort int
	Peers []Peer
}

func (k Key) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}