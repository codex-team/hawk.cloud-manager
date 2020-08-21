package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleConfig(c *gin.Context) {
	rawKey := c.Param("key")
	if rawKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("received empty key")})
		return
	}
	key, err := api.NewKey(rawKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	peers, err := s.matcher.Peers(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	respBody, err := json.Marshal(peers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, respBody)
}

func (s *Server) handleTopology(c *gin.Context) {
	respBody, err := json.Marshal(s.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, respBody)
}
