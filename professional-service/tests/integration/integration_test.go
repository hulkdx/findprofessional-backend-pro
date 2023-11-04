package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, closeDb := InitDb()
	defer closeDb()

	FindAllProfessionalTest(t, db)
	FindAllAvailabilityProfessionalTest(t, db)
	FindAllRatingProfessionalTest(t, db)

	// FindProfessionalTest(t, db, gdb)

	// UpdateProfessionalTest(t, db, gdb)
}
