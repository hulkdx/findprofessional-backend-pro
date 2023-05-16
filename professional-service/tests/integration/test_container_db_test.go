package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDb() (*sql.DB, *gorm.DB, func()) {
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
			WaitingFor: wait.ForListeningPort("5432/tcp"),
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
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = gormDB.AutoMigrate(&professional.Professional{})
	if err != nil {
		log.Fatal(err)
	}

	return db, gormDB, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}
