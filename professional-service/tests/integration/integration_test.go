package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, gdb, closeDb := InitDb()
	defer closeDb()

	FindAllProfessionalTest(t, db, gdb)
	FindProfessionalTest(t, db, gdb)
	FindAllRatingProfessionalTest(t, db, gdb)
	UpdateProfessionalTest(t, db, gdb)
}
