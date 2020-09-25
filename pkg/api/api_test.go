package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	conf = Conf{
		ListenPort: 12345,
		Peers: []Peer{
			{
				PublicKey:                   "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
				PersistentKeepAliveInterval: time.Second,
			},
			{
				PublicKey: "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
			},
		},
	}

	otherConf = Conf{
		ListenPort: 666,
		Peers: []Peer{
			{
				PublicKey:                   "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
				PersistentKeepAliveInterval: time.Second,
			},
			{
				PublicKey: "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
			},
		},
	}
)

// TestEquals tests comparison of Conf objects
func TestEquals(t *testing.T) {
	t.Run("different port", func(t *testing.T) {
		require.False(t, conf.Equals(&otherConf))
	})

	t.Run("different peers", func(t *testing.T) {
		otherConf.ListenPort = conf.ListenPort
		otherConf.Peers[0].Endpoint = "wg.example.com:51820"
		require.False(t, conf.Equals(&otherConf))
	})

	t.Run("different allowed IPs", func(t *testing.T) {
		conf.Peers[0].Endpoint = otherConf.Peers[0].Endpoint
		conf.Peers[0].AllowedIPs = append(conf.Peers[0].AllowedIPs, "10.11.0.76/32", "10.11.0.76/32", "244.12.13.1/32")
		otherConf.Peers[0].AllowedIPs = append(otherConf.Peers[0].AllowedIPs, "10.11.0.76/32", "10.11.0.76/32")
		require.False(t, conf.Equals(&otherConf))
	})

	t.Run("equal", func(t *testing.T) {
		otherConf.Peers[0].AllowedIPs = append(otherConf.Peers[0].AllowedIPs, "244.12.13.1/32")
		require.True(t, conf.Equals(&otherConf))
	})
}

// TestCheckKeyFormat checks public key validation
func TestCheckKeyFormat(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		require.Nil(t, CheckKeyFormat("cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE="))
	})

	t.Run("not base64", func(t *testing.T) {
		require.NotNil(t, CheckKeyFormat("ewerewqwtrtqertqer"))
	})

	t.Run("wrong length", func(t *testing.T) {
		require.NotNil(t, CheckKeyFormat("ZXF3ZXJxd2VydHFlcmV3cmRkZmdoZmdoCg=="))
	})

}
