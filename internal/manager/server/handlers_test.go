package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codex-team/hawk.cloud-manager/internal/manager/storage/yaml"
	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	cfg = config.PeerConfig{
		Hosts: []config.Host{
			{
				Name:       "hawk-collector",
				PublicKey:  "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
				Endpoint:   "10.11.0.2:1234",
				AllowedIPs: []string{"10.11.0.0/24"},
			},
			{
				Name:       "hawk-workers",
				PublicKey:  "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
				Endpoint:   "10.11.0.5:9823",
				AllowedIPs: []string{"10.11.0.0/24"},
			},
		},
		Groups: []config.Group{
			{
				Name:  "hawk-cloud1",
				Hosts: []string{"hawk-collector", "hawk-workers"},
			},
		},
	}

	srv = &Server{
		storage: yaml.NewYamlStorage("cfg.yaml"),
	}

	requestBody = api.Creds{
		PublicKey: "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
	}

	router = srv.setupRouter()
)

// initTests initializes Server fields
func initTests() (err error) {
	srv.storage.Set(cfg)
	return
}

// performRequest mocks HTTP requests
func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

// TestTopology tests getting WireGuard configuration
func TestTopology(t *testing.T) {
	require.Nil(t, initTests())

	// simple case
	t.Run("simple", func(t *testing.T) {
		body, err := json.Marshal(requestBody)
		require.Nil(t, err)

		w := performRequest(router, "POST", "/topology", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)

		wgConf, err := cfg.ToAPIConf(requestBody.PublicKey)
		require.Nil(t, err)

		expected, err := json.Marshal(*wgConf)
		require.Nil(t, err)

		require.Equal(t, string(expected), w.Body.String())
	})

	// sending unknown public key
	t.Run("unknown public key", func(t *testing.T) {
		requestBody.PublicKey = "cnRnaGdmamVkZmdoYndzcmVnd2VyZ2hidHl0cmV5anJx"
		body, err := json.Marshal(requestBody)
		require.Nil(t, err)

		w := performRequest(router, "POST", "/topology", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, `{"failed to get config":"unknown public key"}`, w.Body.String())
	})
}

// TestConfig tests getting Peer Config
func TestConfig(t *testing.T) {
	// simple case
	t.Run("simple", func(t *testing.T) {
		w := performRequest(router, "GET", "/config", nil)
		require.Equal(t, http.StatusOK, w.Code)
		expected, err := json.Marshal(cfg)
		require.Nil(t, err)

		require.Equal(t, string(expected), w.Body.String())
	})
}

// TestUpdate tests updating PeerConfig
func TestUpdateConfig(t *testing.T) {
	// simple case
	t.Run("simple", func(t *testing.T) {
		cfgPatch := cfg
		cfgPatch.Hosts = append(cfgPatch.Hosts, config.Host{
			Name:       "hawk-admin",
			PublicKey:  "cnRnaGdmamVkZmdoYndzcmVnd2VyZ2hidHl0cmV5anJx",
			Endpoint:   "172.17.123.13:3435",
			AllowedIPs: []string{"10.11.0.77/32"},
		})
		body, err := json.Marshal(cfgPatch)
		require.Nil(t, err)
		w := performRequest(router, "PUT", "/config", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, cfgPatch, *srv.storage.Get())
	})
}
