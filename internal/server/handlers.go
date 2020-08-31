package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/gin-gonic/gin"
)

// handleTopology returns WireGuard configuration for Cloud Agent
func (s *Server) handleTopology(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err})
		return
	}
	// Receive public key and signature
	creds := api.Creds{}
	err = json.Unmarshal(bodyBytes, &creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err})
		return
	}

	// Check public key
	for _, h := range s.config.Hosts {
		if h.PublicKey == creds.PublicKey {
			respBody, err := json.Marshal(*s.apiConf)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"failed to write response": err})
				return
			}
			c.JSON(http.StatusOK, respBody)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "unknown public key"})
}
