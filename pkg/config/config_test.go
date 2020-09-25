package config

import (
	"testing"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/stretchr/testify/require"
)

// TestCalendarConfig tests parsing PeerConfig
func TestCalendarConfig(t *testing.T) {
	// simple case
	t.Run("simple", func(t *testing.T) {
		cfg := PeerConfig{
			Hosts: []Host{
				{
					Name:       "hawk-collector",
					PublicKey:  "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
					AllowedIPs: []string{"10.11.0.1/24"},
				},
				{
					Name:       "hawk-workers",
					PublicKey:  "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
					AllowedIPs: []string{"10.11.0.2/24"},
				},
			},
			Groups: []Group{
				{
					Name:  "hawk-cloud1",
					Hosts: []string{"hawk-collector", "hawk-workers"},
				},
			},
		}

		key := "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE="
		err := api.CheckKeyFormat(key)
		require.Nil(t, err)

		expected := api.Conf{
			Peers: []api.Peer{
				{
					PublicKey:  key,
					AllowedIPs: []string{"10.11.0.1/24"},
				},
			},
		}

		apiConf, err := cfg.ToAPIConf("yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=")
		require.Nil(t, err)
		require.Equal(t, expected, *apiConf)
	})
}
