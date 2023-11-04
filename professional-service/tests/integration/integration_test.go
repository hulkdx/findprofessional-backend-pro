package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, closeDb := InitDb()
	defer closeDb()

	FindAllProfessionalTest(t, db)
	// FindProfessionalTest(t, db, gdb)
	FindAllAvailabilityProfessionalTest(t, db)
	// FindAllRatingProfessionalTest(t, db, gdb)
	// UpdateProfessionalTest(t, db, gdb)
}
