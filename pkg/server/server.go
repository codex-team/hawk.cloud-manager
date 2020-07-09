package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type server struct {
	config  config.PeerConfig
	port    string
	logger  zap.Logger
	Handler http.Handler
}

func (s *server) Run() {
	fmt.Println("Server")
}

func HandleConfig(cfg *config.PeerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func New(port string, config config.PeerConfig, logger zap.Logger) *server {
	router := mux.NewRouter()
	router.HandleFunc("/config", HandleConfig(&config))

	var server = &server{
		port:    port,
		config:  config,
		logger:  logger,
		Handler: router,
	}
	return server
}
