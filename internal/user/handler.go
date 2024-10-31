package user

import "github.com/gin-gonic/gin"


type UserHandler interface {
	RegisterUser(c *gin.Context)
	GetUserByID(c *gin.Context)
}