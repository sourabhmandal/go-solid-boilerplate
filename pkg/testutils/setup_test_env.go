package testutils

import (
	"authosaur/pkg/database"
	"authosaur/web"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestServer struct {
	router           *gin.Engine
	httpServer     web.WebServer
	dbContainer testcontainers.Container
	dbInst 		database.Database
	db 				*pgx.Conn
}

var TestServerInst *TestServer

func StartTestServer(ctx context.Context) *TestServer {
	if(TestServerInst != nil) {
		return TestServerInst
	}

	dbContainer, dbInst, db := TestServerInst.startPostgresContainer(ctx)
	// Initialize the Gin router and the HTTP server
	httpServer := web.NewHttpServer(8080, dbInst, db)
	httpServer.Start()
	router := httpServer.GetRouter()

	return &TestServer{
		router: router,
		httpServer: httpServer,
		dbContainer: dbContainer,
		dbInst: dbInst,
		db: db,
	}
}

func (ts *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts.router.ServeHTTP(w, r)
}

func(ts *TestServer) StopTestServer(ctx context.Context) {
	// cleanups
	log.Println("Cleanup test server with postgres, gin router and http-server")

	if err := ts.dbContainer.Terminate(ctx); err != nil {
    log.Println(err)
  }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ts.dbInst.Close(); err != nil {
		log.Printf("Database unable to stop with error: %v", err)
	}
	if err := ts.httpServer.GetServer().Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
}

// startPostgresContainer initializes the PostgreSQL test container.
func(ts *TestServer) startPostgresContainer(ctx context.Context) (testcontainers.Container, database.Database, *pgx.Conn) {
  var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

  dbContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Println(err)
    return nil, nil, nil
	}
	dbHost, err := dbContainer.Host(ctx)
	if err != nil {
		log.Println(err)
    return nil, nil, nil
	}

	dbPort, err := dbContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
    log.Println(err)
    return nil, nil, nil
	}

  
	if err := dbContainer.Start(ctx); err != nil {
		log.Fatalf("Failed to start container: %v", err)
    return nil, nil, nil
	}

	host := dbHost
	port, _ := strconv.Atoi(dbPort.Port())
  dbInst, db := database.NewDatabasePg(dbUser, dbPwd, host, dbName, "public", port)
	return dbContainer, dbInst, db
}
