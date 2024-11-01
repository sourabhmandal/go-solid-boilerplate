package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/net/context"
)

// Service represents a service that interacts with a database.
type Database interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type pgDatabase struct {
	db *pgx.Conn
}

var (
	dbInstance *pgDatabase
)

func NewDatabasePg(username, password, host, database, schema string, port int) (Database, *pgx.Conn) {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, dbInstance.db
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &pgDatabase{
		db: db,
	}
	return dbInstance, db
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *pgDatabase) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *pgDatabase) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close(context.Background())
}
