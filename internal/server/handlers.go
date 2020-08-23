package server

import (
	"io/ioutil"
	"net/http"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
	"github.com/gin-gonic/gin"
)

/*
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
*/

func (s *Server) handleTopology(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to read request body": err})
		return
	}
	creds := api.Creds{}
	err = creds.UnmarshalJSON(bodyBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"failed to parse JSON": err})
		return
	}

	for _, h := range s.config.Hosts {
		if h.PublicKey == creds.PublicKey {
			respBody, err := (*s.apiConf).MarshalJSON()
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
