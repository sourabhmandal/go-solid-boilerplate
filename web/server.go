package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebServer interface {
  Start() error
  GracefulShutdown(done chan bool)
  GetRouter() *gin.Engine
  GetServer() *http.Server
}