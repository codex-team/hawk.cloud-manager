package server

import (
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config  *config.PeerConfig
	apiConf *api.Conf
	http    *http.Server
}

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

func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/topology", s.handleTopology)

	return r
}

func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)
	return s.http.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.http.Close()
}
