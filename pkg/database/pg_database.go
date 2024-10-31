package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

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

	// Get database stats (like open connections, in use, idle, etc.)
	// dbStats := s.db.Stat()
	// stats["open_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	// stats["idle"] = strconv.Itoa(int(dbStats.IdleConns()))
	// stats["wait_count"] = strconv.FormatInt(int64(dbStats.ConstructingConns()), 10)

	// // Evaluate stats to provide a health message
	// if dbStats.AcquiredConns() > 40 { // Assuming 50 is the max for this example
	// 	stats["message"] = "The database is experiencing heavy load."
	// }

	// if dbStats.MaxIdleDestroyCount() > int64(dbStats.AcquireCount())/2 {
	// 	stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	// }

	// if dbStats.MaxLifetimeDestroyCount() > int64(dbStats.AcquireCount())/2 {
	// 	stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	// }

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
