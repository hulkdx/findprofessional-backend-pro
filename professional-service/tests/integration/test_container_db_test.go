package integration_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func InitDb(t *testing.T) (*pgxpool.Pool, func()) {
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
		t.Fatal("Failed to start PostgreSQL container: ", err)
	}

	// Get the container's host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal("Failed to get PostgreSQL container host: ", err)
	}
	port, err := postgresContainer.MappedPort(ctx, nat.Port("5432"))
	if err != nil {
		t.Fatal("Failed to get PostgreSQL container port: ", err)
	}

	// Construct the PostgreSQL connection string
	connectionString := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", host, port.Port())

	// Connect to the PostgreSQL container
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		t.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	err = integrationDbTables(ctx, pool)
	if err != nil {
		t.Fatal(err)
	}

	return pool, func() {
		pool.Close()
		postgresContainer.Terminate(ctx)
	}
}

func integrationDbTables(ctx context.Context, db *pgxpool.Pool) error {
	content, err := fetchURLContent("https://raw.githubusercontent.com/hulkdx/findprofessional-backend-user/main/user-service/src/main/resources/db/changelog/db.changelog-master.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, content)
	return err
}

func fetchURLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	content := string(body)
	return content, nil
}
