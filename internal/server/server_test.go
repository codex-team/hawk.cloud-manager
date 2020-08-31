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
		config: &cfg,
	}

	requestBody = api.Creds{
		PublicKey: "cXdlcnRydGV3cnd0cnRxcnFlcnFydHRydHJ5dXlyZXE=",
		Signature: "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk=",
	}
)

func initTests() (err error) {
	srv.apiConf, err = cfg.ToAPIConf()
	return
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestTopology(t *testing.T) {
	require.Nil(t, initTests())
	t.Run("simple", func(t *testing.T) {
		body, err := json.Marshal(requestBody)
		require.Nil(t, err)

		router := srv.setupRouter()
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

	t.Run("unknown public key", func(t *testing.T) {
		requestBody.PublicKey = "yAnz5TF+lXXJte14tji3zlMNq+hd2rYUIgJBgB3fBmk="
		body, err := json.Marshal(requestBody)
		require.Nil(t, err)

		router := srv.setupRouter()
		w := performRequest(router, "POST", "/topology", bytes.NewReader(body))
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, `{"error":"unknown public key"}`, w.Body.String())
	})
}
