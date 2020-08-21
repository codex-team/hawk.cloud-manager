package server

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/akath19/gin-zap"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/codex-team/hawk.cloud-manager/pkg/matcher"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const logDuration = 3 * time.Second

type Server struct {
	config  *config.PeerConfig
	matcher *matcher.Simple
	http    *http.Server
}

func New(addr string, config *config.PeerConfig, logger *zap.Logger) (*Server, error) {
	peerMatcher, err := matcher.NewSimpleMatcher(*config)
	if err != nil {
		return nil, err
	}

	var server = &Server{
		config:  config,
		matcher: peerMatcher,
		http: &http.Server{
			Addr: addr,
		},
	}
	server.http.Handler = server.setupRouter(logger)

	return server, nil
}

func (s *Server) setupRouter(logger *zap.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(ginzap.Logger(logDuration, logger))

	r.GET("/topology", s.handleTopology)
	r.GET("/config/:key", s.handleConfig)

	return r
}

func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)
	fmt.Println("Server")

	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}
