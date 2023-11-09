package integration_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
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
	(id,"email","password","first_name","last_name","coach_type","price_number","price_currency", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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
			p.CreatedAt,
			p.UpdatedAt,
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
	query := `INSERT INTO "professional_availability" 
	("professional_id", "date", "from", "to", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5, $6)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, a := range availabilities {
		_, err := stmt.Exec(
			a.ProfessionalID,
			a.Date.String(),
			a.From.String(),
			a.To.String(),
			a.CreatedAt,
			a.UpdatedAt,
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
		defer db.Exec(`DELETE FROM professional_availability; DELETE FROM professionals;`)
	}
}

func insertReview(db *sql.DB, review ...professional.Review) func() {
	query := `INSERT INTO "professional_review" 
	("professional_id", "user_id", "rate", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, r := range review {
		_, err := stmt.Exec(
			r.ProfessionalID,
			r.UserID,
			r.Rate,
			r.CreatedAt,
			r.UpdatedAt,
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
		defer db.Exec(`DELETE FROM professional_review;`)
	}
}

func insertUserWithId(db *sql.DB, userId ...int) func() {
	u := []user.User{}
	for _, id := range userId {
		u = append(u, user.User{ID: id, Email: fmt.Sprint(id)})
	}
	return insertUser(db, u...)
}

func insertUser(db *sql.DB, user ...user.User) func() {
	query := `INSERT INTO "users"
	(id, email, password, first_name, last_name, profile_image) VALUES
	($1, $2, $3, '', $4, $5)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, u := range user {
		_, err := stmt.Exec(
			u.ID,
			u.Email,
			u.FirstName,
			u.LastName,
			u.ProfileImage,
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
		defer db.Exec(`DELETE FROM users;`)
	}
}
