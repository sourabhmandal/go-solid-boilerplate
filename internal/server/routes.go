package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterServerRoutes() http.Handler {
	r := gin.New()
	r.GET("/", s.helloWorldHandler)
	r.GET("/health", s.healthHandler)

	return r
}

func (s *Server) helloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
