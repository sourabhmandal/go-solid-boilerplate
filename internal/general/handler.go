package general

import "github.com/gin-gonic/gin"


type GeneralHandler interface {
	HealthCheck(c *gin.Context)
}