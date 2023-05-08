package config

import (
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func TestDatabaseConfigConnectionString(t *testing.T) {
	// Arrange
	expected := "postgresql://postgres:postgres@postgresdb/postgres"
	url := "postgresql://postgresdb/postgres"
	username := "postgres"
	password := "postgres"
	db := DatabaseConfig{url: url, username: username, password: password}
	// Act
	result := db.ConnectionString()
	// Assert
	assert.Equal(t, result, expected)
}
