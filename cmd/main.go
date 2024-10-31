package main

import (
	"authosaurous/pkg/database"
	"authosaurous/web"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)



func main() {
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	databaseName := os.Getenv("DB_DATABASE")
	password   := os.Getenv("DB_PASSWORD")
	username   := os.Getenv("DB_USERNAME")
	host       := os.Getenv("DB_HOST")
	dbport, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	schema     := os.Getenv("DB_SCHEMA")

	dbInst, db := database.NewDatabasePg(username, password, host, databaseName, schema, dbport)
	newServer := web.NewHttpServer(serverPort, dbInst, db)
	// Run graceful shutdown in a separate goroutine
	go newServer.GracefulShutdown(done)

	if err:= newServer.Start();err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
