package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/codex-team/hawk.cloud-manager/pkg/matcher"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestServer(t *testing.T) {
	cfg := config.PeerConfig{
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
	m, err := matcher.NewSimpleMatcher(cfg)
	require.Nil(t, err)
	srv := &Server{
		config:  &cfg,
		matcher: m,
	}
	logger, err := zap.NewDevelopment()
	require.Nil(t, err)

	t.Run("config", func(t *testing.T) {
		router := srv.setupRouter(logger)
		w := performRequest(router, "GET", "/config/cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=", nil)
		require.Equal(t, http.StatusOK, w.Code)
		strBody, err := strconv.Unquote(w.Body.String())
		require.Nil(t, err)
		data, err := base64.StdEncoding.DecodeString(strBody)
		require.Nil(t, err)
		peers := make([]api.Peer, 0)
		err = json.Unmarshal(data, &peers)
		require.Nil(t, err)
	})

	t.Run("topology", func(t *testing.T) {
		router := srv.setupRouter(logger)
		w := performRequest(router, "GET", "/topology", nil)
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
