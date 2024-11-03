package general

import (
	"gosolid/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeneralHandler interface {
	HealthCheck(c *gin.Context)
}

type generalHandler struct {
	db database.Database
}

func NewGeneralHandler(db database.Database) GeneralHandler {
	return &generalHandler{
		db: db,
	}
}

func (gh *generalHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"server":   "ok",
		"database": gh.db.Health(),
	})
}
