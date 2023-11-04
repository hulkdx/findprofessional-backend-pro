package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"cloud.google.com/go/civil"
	"github.com/docker/go-connections/nat"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}),
		&gorm.Config{
			NamingStrategy: TestNamingStrategy{},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	integrationDbTables(gormDB)
	if err != nil {
		log.Fatal(err)
	}

	return db, gormDB, func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}
}

func integrationDbTables(gormDB *gorm.DB) error {
	schema.RegisterSerializer("custom", CustomSerializer{})
	err := gormDB.AutoMigrate(&professional.Professional{}, &professional.Availability{})
	if err != nil {
		return err
	}
	// https://github.com/hulkdx/findprofessional-backend-user/blob/main/user-service/src/main/resources/db/changelog/db.changelog-master.sql
	err = gormDB.AutoMigrate(&ProfessionalRating{})
	if err != nil {
		return err
	}
	return nil
}

type TestNamingStrategy struct {
	schema.NamingStrategy
}

func (ns TestNamingStrategy) TableName(table string) string {
	if table == "ProfessionalRating" {
		return "professional_rating"
	}
	if table == "Availability" {
		return "professional_availability"
	}
	return ns.NamingStrategy.TableName(table)
}

type ProfessionalRating struct {
	ID             uint `gorm:"primaryKey"`
	UserID         int64
	ProfessionalID int64
	Rate           int
}

type CustomSerializer struct{}

func (CustomSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	t := sql.NullTime{}
	if err = t.Scan(dbValue); err == nil && t.Valid {
		err = field.Set(ctx, dst, t.Time.Unix())
	}
	return
}

func (CustomSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (result interface{}, err error) {
	switch v := fieldValue.(type) {
	case civil.Time:
		result = v.String()
	default:
		err = fmt.Errorf("invalid field type %#v for UnixSecondSerializer, only int, uint supported", v)
	}
	return
}
