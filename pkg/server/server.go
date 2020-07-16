package server

import (
	"encoding/json"
	"fmt"
	"github.com/codex-team/hawk.cloud-manager/pkg/matcher"
	"github.com/codex-team/hawk.cloud-manager/pkg/utils"
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

func HandleTopology(cfg *config.PeerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleConfig(m matcher.Matcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rawKey := vars["key"]
		key, err := utils.ParseKey(rawKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		peers,err := m.Peers(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(peers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func New(port string, config config.PeerConfig, logger zap.Logger) (*server, error) {
	peerMatcher, err := matcher.NewSimpleMatcher(config)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	router.HandleFunc("/topology", HandleTopology(&config))
	router.HandleFunc("/config/{key}", HandleConfig(peerMatcher))

	var server = &server{
		port:    port,
		config:  config,
		logger:  logger,
		Handler: router,
	}
	return server, nil
}
