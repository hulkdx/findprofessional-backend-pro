package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, gdb, closeDb := InitDb()
	defer closeDb()

	t.Run("FindAllProfessionalTests", func(t *testing.T) {
		FindAllProfessionalTest(t, db, gdb)
	})
	t.Run("FindProfessionalTests", func(t *testing.T) {
		FindProfessionalTest(t, db, gdb)
	})
}
