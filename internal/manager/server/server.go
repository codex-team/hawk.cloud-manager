package server

import (
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/internal/manager/storage"
	"github.com/gin-gonic/gin"
)

// Server is struct that wraps HTTP server and contains current configuration
type Server struct {
	// Storage of configs
	storage storage.Storage
	// HTTP server
	http *http.Server
}

// New creates Server and returns pointer to it
func New(addr string, st storage.Storage) (*Server, error) {
	var server = &Server{
		storage: st,
		http: &http.Server{
			Addr: addr,
		},
	}
	server.http.Handler = server.setupRouter()

	return server, nil
}

// setupRouter initializes handlers and returns new router
func (s *Server) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/topology", s.handleTopology)
	r.GET("/config", s.handleConfig)
	r.PUT("/config", s.handleUpdate)

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
