package config

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
	"testing"
)

func TestDatabaseConfigConnectionString(t *testing.T) {
	t.Run("connection string should add username password", func(t *testing.T) {
		// Arrange
		expected := "postgresql://postgres:postgres@postgresdb/postgres?sslmode=disable"
		url := "postgresql://postgresdb/postgres"
		username := "postgres"
		password := "postgres"
		db := DatabaseConfig{url: url, username: username, password: password}
		// Act
		result := db.ConnectionString()
		// Assert
		assert.Equal(t, result, expected)
	})

	t.Run("connection string with sslmode", func(t *testing.T) {
		// Arrange
		expected := "postgresql://postgres:postgres@postgresdb/postgres?sslmode=require"
		url := "postgresql://postgresdb/postgres?sslmode=require"
		username := "postgres"
		password := "postgres"
		db := DatabaseConfig{url: url, username: username, password: password}
		// Act
		result := db.ConnectionString()
		// Assert
		assert.Equal(t, result, expected)
	})
}
