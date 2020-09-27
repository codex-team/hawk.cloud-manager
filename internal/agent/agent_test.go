package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/stretchr/testify/require"
)

var (
	conf = api.Conf{
		ListenPort: 12345,
		Peers: []api.Peer{
			{
				PublicKey:                   "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
				Endpoint:                    "10.11.12.1:5454",
				PersistentKeepAliveInterval: 5 * time.Second,
				AllowedIPs:                  []string{"234.12.122.0/32"},
			},
			{
				PublicKey:  "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
				Endpoint:   "234.12.122.54:7823",
				AllowedIPs: []string{"10.11.12.0/16"},
			},
		},
	}
)

// TestQueryConfig tests receiving WireGuard config from Manager
func TestQueryConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, req.URL.String(), "/topology")
			respBody, err := json.Marshal(conf)
			require.Nil(t, err)
			_, err = rw.Write(respBody)
			require.Nil(t, err)
		}))
		defer server.Close()

		agent := Agent{
			client:         server.Client(),
			managerAddress: server.URL,
			PubKey:         "cnRnaGdmamVkZmdoYndzcmVnd2VyZ2hidHl0cmV5anJx",
		}
		returnedConf, err := agent.queryConf()
		require.Nil(t, err)
		require.True(t, returnedConf.Equals(&conf))
	})
}

// TestParseConfig tests parsing api.Conf to WireGuard config
func TestParseConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		agent := Agent{
			PrivKey: "HIgo9xNzJMWLKASShiTqIybxZ0U3wGLiUeJ1PKf8ykw=",
			PubKey:  "cnRnaGdmamVkZmdoYndzcmVnd2VyZ2hidHl0cmV5anJx",
			Config:  &conf,
		}
		res, err := agent.parseWGConf()
		require.Nil(t, err)

		var expected = `[Interface]
PrivateKey = HIgo9xNzJMWLKASShiTqIybxZ0U3wGLiUeJ1PKf8ykw=
ListenPort = 12345

[Peer]
PublicKey = cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=
Endpoint = 10.11.12.1:5454
AllowedIPs = 234.12.122.0/32

[Peer]
PublicKey = yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=
Endpoint = 234.12.122.54:7823
AllowedIPs = 10.11.12.0/16`

		require.Equal(t, expected, res)
	})
}
