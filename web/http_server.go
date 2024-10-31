package web

import (
	"authosaur/internal/general"
	"authosaur/internal/repository"
	"authosaur/internal/user"
	"authosaur/pkg/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type HTTPServer struct {
	port int
	server *http.Server
	router *gin.Engine
	db database.Database
}

func NewHttpServer(serverPort int, dbInst database.Database, db *pgx.Conn) WebServer {
	// Declare Router
	queries := repository.New(db)
	// declare handlers
	generalHandler := general.NewGeneralHandler(dbInst)
	userService := user.NewUserService(queries)
	userHandlers := user.NewUserHandler(userService)

	// declare routes
	router := gin.Default()
	// generic routes
	router.GET("/health", generalHandler.HealthCheck)
	// user routes
	userRouter := router.Group("/user")
	userRouter.POST("/", userHandlers.RegisterUser)
	userRouter.GET("/:id", userHandlers.GetUserByID)
	
	newServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: router,
	}
	
	return &HTTPServer{
		port: serverPort,
		server: newServer,
		router: router,
		db: dbInst,
	}
}

func (s *HTTPServer) GetRouter() *gin.Engine {
	return s.router
}

func (s *HTTPServer) GetServer() *http.Server {
	return s.server
}


func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) GracefulShutdown(done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.db.Close(); err != nil {
		log.Printf("Database unable to stop with error: %v", err)
	}
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}