package integration_test

import (
	"testing"
)

func TestIntegrations(t *testing.T) {
	db, closeDb := InitDb(t)
	defer closeDb()

	t.Run("FindAllProfessionalTest", func(t *testing.T) {
		FindAllProfessionalTest(t, db)
	})
	t.Run("FindAllAvailabilityProfessionalTest", func(t *testing.T) {
		FindAllAvailabilityProfessionalTest(t, db)
	})
	t.Run("FindAllRatingProfessionalTest", func(t *testing.T) {
		FindAllRatingProfessionalTest(t, db)
	})
	t.Run("FindAllReviewProfessionalTest", func(t *testing.T) {
		FindAllReviewProfessionalTest(t, db)
	})

	t.Run("FindProfessionalTest", func(t *testing.T) {
		FindProfessionalTest(t, db)
	})

	t.Run("UpdateProfessionalTest", func(t *testing.T) {
		UpdateProfessionalTest(t, db)
	})

	t.Run("ReviewFindAllTest", func(t *testing.T) {
		ReviewFindAllTest(t, db)
	})

	t.Run("CreateProTest", func(t *testing.T) {
		CreateProTest(t, db)
	})

	t.Run("GetAvailabilityTest", func(t *testing.T) {
		GetAvailabilityTest(t, db)
	})

	t.Run("AddAvailabilityTest", func(t *testing.T) {
		AddAvailabilityTest(t, db)
	})
}
