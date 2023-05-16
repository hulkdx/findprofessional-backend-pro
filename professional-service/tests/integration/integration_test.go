package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, gdb, closeDb := InitDb()
	defer closeDb()

	t.Run("ListProfessionalTests", func(t *testing.T) {
		ListProfessionalTest(t, db, gdb)
	})
}
