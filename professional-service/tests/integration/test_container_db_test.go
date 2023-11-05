package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	integrationDbTables(db)
	if err != nil {
		log.Fatal(err)
	}

	return db, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}

func integrationDbTables(db *sql.DB) error {
	// https://github.com/hulkdx/findprofessional-backend-user/blob/main/user-service/src/main/resources/db/changelog/db.changelog-master.sql
	_, err := db.Exec(`
	CREATE TABLE "users" (
		"id"  BIGSERIAL PRIMARY KEY,
		"email" VARCHAR(255) UNIQUE NOT NULL,
		"password" VARCHAR(255) NOT NULL,
		"first_name" VARCHAR(255),
		"last_name" VARCHAR(255),
		"created_at" timestamp,
		"updated_at" timestamp
	);
	
	CREATE TABLE "professionals" (
		"id" BIGSERIAL PRIMARY KEY,
		"email" VARCHAR(255) UNIQUE NOT NULL,
		"password" VARCHAR(255) NOT NULL,
		"first_name" VARCHAR(255) NOT NULL,
		"last_name" VARCHAR(255) NOT NULL,
		"coach_type" VARCHAR(255) NOT NULL,
		"price_number" BIGINT NOT NULL,
		"price_currency" VARCHAR(255) NOT NULL,
		"profile_image_url" VARCHAR(255),
		"description" VARCHAR(255),
		"created_at" timestamp,
		"updated_at" timestamp
	);
	
	CREATE TABLE "professional_rating" (
		"id" BIGSERIAL PRIMARY KEY,
		"user_id" BIGINT NOT NULL,
		"professional_id" BIGINT NOT NULL,
		"rate" INT NOT NULL
	);
	
	CREATE UNIQUE INDEX ON "professional_rating" ("user_id", "professional_id");
	
	ALTER TABLE "professional_rating" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
	
	ALTER TABLE "professional_rating" ADD FOREIGN KEY ("professional_id") REFERENCES "professionals" ("id");
	
	CREATE TABLE "professional_availability" (
		"id" BIGSERIAL PRIMARY KEY,
		"professional_id" BIGINT NOT NULL,
		"date" DATE NOT NULL,
		"from" TIME NOT NULL,
		"to" TIME NOT NULL
	);
	
	ALTER TABLE "professional_availability" ADD FOREIGN KEY ("professional_id") REFERENCES "professionals" ("id");
	`)
	return err
}
