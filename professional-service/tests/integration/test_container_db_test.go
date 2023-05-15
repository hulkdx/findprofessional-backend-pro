package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"

	_ "github.com/lib/pq"
)

func InitDb() (*sql.DB, func()) {
	ctx := context.Background()

	// Create a PostgreSQL container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpassword",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal("Failed to start PostgreSQL container: ", err)
	}

	// Get the container's host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatal("Failed to get PostgreSQL container host: ", err)
	}
	port, err := postgresContainer.MappedPort(ctx, nat.Port("5432"))
	if err != nil {
		log.Fatal("Failed to get PostgreSQL container port: ", err)
	}

	// Construct the PostgreSQL connection string
	connectionString := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", host, port.Port())

	// Connect to the PostgreSQL container
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL: ", err)
	}

	return db, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}
