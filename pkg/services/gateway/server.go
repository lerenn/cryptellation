//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ../../../api/gateway/v1.yaml

package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lerenn/cryptellation/v1/clients/go/client"
)

// Server is the gateway server.
type Server struct {
	client client.Client
}

// NewServer creates a new gateway server.
func NewServer(client client.Client) *Server {
	return &Server{
		client: client,
	}
}

// GetInfo is the handler for the GET /info endpoint.
func (s *Server) GetInfo(c *gin.Context) {
	info, err := s.client.Info(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SystemInformation{
		Version: &info.Version,
	})
}
