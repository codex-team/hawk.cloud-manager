package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/codex-team/hawk.cloud-manager/internal/storage/yaml"
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
				AllowedIPs: []string{"10.11.0.1/24"},
			},
			{
				Name:       "hawk-workers",
				PublicKey:  "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
				AllowedIPs: []string{"10.11.0.2/24"},
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
		Signature: "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
	}

	router = srv.setupRouter()
)

// initTests initializes Server fields
func initTests() (err error) {
	srv.apiConf, err = cfg.ToAPIConf()
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

		expected, err := json.Marshal(srv.apiConf)
		require.Nil(t, err)

		strBody, err := strconv.Unquote(w.Body.String())
		require.Nil(t, err)
		actual, err := base64.StdEncoding.DecodeString(strBody)
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})

	// sending unknown public key
	t.Run("unknown public key", func(t *testing.T) {
		requestBody.PublicKey = "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk="
		body, err := json.Marshal(requestBody)
		require.Nil(t, err)

		w := performRequest(router, "POST", "/topology", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, `{"error":"unknown public key"}`, w.Body.String())
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
		strBody, err := strconv.Unquote(w.Body.String())
		require.Nil(t, err)
		actual, err := base64.StdEncoding.DecodeString(strBody)
		require.Nil(t, err)

		require.Equal(t, expected, actual)
	})
}

// TestUpdate tests updating PeerConfig
func TestUpdateConfig(t *testing.T) {
	// simple case
	t.Run("simple", func(t *testing.T) {
		cfgPatch := cfg
		cfgPatch.Hosts = append(cfgPatch.Hosts, config.Host{
			Name:       "hawk-admin",
			PublicKey:  "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
			Endpoint:   "172.17.123.13",
			AllowedIPs: []string{"10.11.0.77/32"},
		})
		body, err := json.Marshal(cfgPatch)
		require.Nil(t, err)
		w := performRequest(router, "PUT", "/config", bytes.NewReader(body))
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, cfgPatch, *srv.storage.Get())
	})
}
