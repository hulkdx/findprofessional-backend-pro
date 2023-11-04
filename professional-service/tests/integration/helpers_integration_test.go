package integration_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http/httptest"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func Unmarshal(response *httptest.ResponseRecorder, output any) {
	json.Unmarshal(response.Body.Bytes(), output)
}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

// Database helpers

func insertEmptyPro(db *sql.DB) func() {
	return insertPro(db, professional.Professional{})
}

func insertPro(db *sql.DB, pro ...professional.Professional) func() {

	query := `INSERT INTO "professionals"
	(id,"email","password","first_name","last_name","coach_type","price_number","price_currency") VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, p := range pro {
		if p.PriceNumber == nil {
			p.PriceNumber = Int(0)
		}
		if p.PriceCurrency == nil {
			p.PriceCurrency = String("")
		}
		_, err := stmt.Exec(
			p.ID,
			p.Email,
			p.Password,
			p.FirstName,
			p.LastName,
			p.CoachType,
			p.PriceNumber,
			p.PriceCurrency,
		)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return func() {
		db.Exec(`DELETE FROM professionals`)
	}
}

func insertAvailability(db *sql.DB, availabilities ...professional.Availability) func() {
	query := `INSERT INTO "professional_availability" ("professional_id", "date", "from", "to") VALUES ($1, $2, $3, $4)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, avail := range availabilities {
		_, err := stmt.Exec(avail.ProfessionalID, avail.Date.String(), avail.From.String(), avail.To.String())
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	return func() {
		defer db.Exec(`DELETE FROM professional_availability; DELETE FROM professionals;`)
	}
}
