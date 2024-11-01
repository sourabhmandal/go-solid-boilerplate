package main

import (
	"authosaur/internal/general"
	"authosaur/internal/repository"
	"authosaur/internal/user"
	"authosaur/pkg/database"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)



func main() {
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	databaseName := os.Getenv("DB_DATABASE")
	password   := os.Getenv("DB_PASSWORD")
	username   := os.Getenv("DB_USERNAME")
	host       := os.Getenv("DB_HOST")
	dbport, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	schema     := os.Getenv("DB_SCHEMA")

	// Declare a flag to run migrations only
  migrateFlag := flag.Bool("migrate", false, "run migrations only")
	// Parse the flags
	flag.Parse()

	if *migrateFlag {
		migrateDatabase(username, password, host, databaseName, schema, dbport)
		return
	}
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	dbInst, db := database.NewDatabasePg(username, password, host, databaseName, schema, dbport)
	
	newServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: registerRoutes(dbInst, db),
	}
	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(done, newServer, dbInst)

	// start the server
	if err:= newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}


func registerRoutes(dbInst database.Database, db *pgx.Conn) *gin.Engine {
	// Declare Router
	queries := repository.New(db)

	// declare generic handlers
	generalHandlers := general.NewGeneralHandler(dbInst)
	// declare user handlers
	userService := user.NewUserService(queries)
	userHandlers := user.NewUserHandler(userService)
	
	router:= gin.Default()
	// generic routes
	router.GET("/health", generalHandlers.HealthCheck)
	// user routes
	userRouter := router.Group("/user")
	userRouter.POST("/", userHandlers.RegisterUser)
	userRouter.GET("/", userHandlers.GetAllUsers)

	return router
}

func gracefulShutdown(done chan bool, server *http.Server, db database.Database) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// shut down any database connections
	if err := db.Close(); err != nil {
		log.Printf("Database unable to stop with error: %v", err)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}


func migrateDatabase(username, password, host, databaseName, schema string, dbport int) {
		// Construct the connection string
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", 
			username, password, host, dbport, databaseName, schema)
	
		// Create a new migrate instance
		m, err := migrate.New("file://pkg/schema", connStr)
		if err != nil {
			log.Fatalf("failed to create migrate instance: %v", err)
		}
	
		// Apply migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	
		log.Println("Migrations applied successfully.")
}
