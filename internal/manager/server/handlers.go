package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/gin-gonic/gin"
)

// handleTopology returns WireGuard configuration for Cloud Agent
func (s *Server) handleTopology(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}
	// Receive public key and signature
	creds := api.Creds{}
	err = json.Unmarshal(bodyBytes, &creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}

	conf, err := (*s.storage.Get()).ToAPIConf(creds.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to get config": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *conf)
}

// handleConfig returns current Peer Config
func (s *Server) handleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, *s.storage.Get())
}

// handleUpdate updates Peer Config with received data
func (s *Server) handleUpdate(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err.Error()})
		return
	}
	// Receive updated version
	cfgPatch := config.PeerConfig{}
	err = json.Unmarshal(bodyBytes, &cfgPatch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err.Error()})
		return
	}
	// Update Peer Config
	s.storage.Set(cfgPatch)
	err = s.storage.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"failed to save config": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
