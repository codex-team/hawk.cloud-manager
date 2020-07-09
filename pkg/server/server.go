package server

import (
	"fmt"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
)

type server struct {
	config config.PeerConfig
	port string
}

func (s *server) Run() {
	fmt.Println("Server")
}

func NewServer(port string, config config.PeerConfig) *server {
	var server = &server{
		port: port,
		config: config,
	}
	return server
}
