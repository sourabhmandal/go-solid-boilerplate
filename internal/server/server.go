package server

import (
	"authosaurus/pkg/database"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supareel/authosaurous/internal/user"
	"github.com/supareel/authosaurous/pkg/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *Server {

	// Declare Router
	router := gin.Default()
	db := database.NewService()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	// user handlers
	userQueries := user.NewUserUseCase()
	userUseCase := user.NewUserUseCase(userQueries)
	userHandler := user.NewUserHandler(userUseCase)



	newServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}




	
	return new 
}
