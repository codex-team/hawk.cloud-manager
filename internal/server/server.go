package server

import (
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/gin-gonic/gin"
)

// Server is struct that wraps HTTP server and contains current configuration
type Server struct {
	// Peer config
	config *config.PeerConfig
	// WireGuard config
	apiConf *api.Conf
	// HTTP server
	http *http.Server
}

// New creates Server and returns pointer to it
func New(addr string, config *config.PeerConfig) (*Server, error) {
	var server = &Server{
		config: config,
		http: &http.Server{
			Addr: addr,
		},
	}
	server.http.Handler = server.setupRouter()
	apiConf, err := (*server.config).ToAPIConf()
	if err != nil {
		return nil, err
	}
	server.apiConf = apiConf

	return server, nil
}

// setupRouter initializes handlers and returns new router
func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/topology", s.handleTopology)

	return r
}

// Run starts HTTP server
func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

// Stop closes HTTP server
func (s *Server) Stop() error {
	return s.http.Close()
}
