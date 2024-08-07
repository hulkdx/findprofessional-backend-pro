package integration_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func OutputSQL(t *testing.T, db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		t.Fatal(err)
	}

	values := make([]sql.RawBytes, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	fmt.Println(strings.Join(columns, "\t\t")) // Print column names

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			t.Fatal(err)
		}

		var rowStrings []string
		for _, raw := range values {
			if raw == nil {
				rowStrings = append(rowStrings, "NULL")
			} else {
				rowStrings = append(rowStrings, string(raw))
			}
		}
		rowString := strings.Join(rowStrings, "\t\t")
		fmt.Println(rowString)
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func insertEmptyPro(t *testing.T, db *sql.DB) func() {
	return insertPro(t, db, professional.Professional{
		PriceNumber:   Int(0),
		PriceCurrency: String(""),
		Pending:       false,
	})
}

func insertPro(t *testing.T, db *sql.DB, pro ...professional.Professional) func() {

	query := `INSERT INTO "professionals"
	(id,"email","password","first_name","last_name","coach_type","price_number","price_currency", "pending", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	for _, p := range pro {
		_, err := stmt.Exec(
			p.ID,
			p.Email,
			p.Password,
			p.FirstName,
			p.LastName,
			p.CoachType,
			p.PriceNumber,
			p.PriceCurrency,
			p.Pending,
			p.CreatedAt,
			p.UpdatedAt,
		)
		if err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	return func() {
		db.Exec(`DELETE FROM professionals`)
	}
}

func insertAvailability(t *testing.T, db *sql.DB, availabilities ...professional.Availability) func() {
	query := `INSERT INTO "professional_availability" 
	("professional_id", "date", "from", "to", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5, $6)`

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		t.Fatal(err)
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
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	return func() {
		defer db.Exec(`DELETE FROM professional_availability; DELETE FROM professionals;`)
	}
}

func insertReview(t *testing.T, db *sql.DB, review ...professional.Review) func() {
	query := `INSERT INTO "professional_review" 
	("professional_id", "user_id", "rate", "created_at", "updated_at", "content_text", "id") VALUES
	($1, $2, $3, $4, $5, $6, $7)`

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range review {
		_, err := db.Exec(
			query,
			r.ProfessionalID,
			r.UserID,
			r.Rate,
			r.CreatedAt,
			r.UpdatedAt,
			r.ContentText,
			r.ID,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	return func() {
		defer db.Exec(`DELETE FROM professional_review;`)
	}
}

func insertUserWithId(t *testing.T, db *sql.DB, userId ...int) func() {
	u := []user.User{}
	for _, id := range userId {
		u = append(u, user.User{ID: id, Email: fmt.Sprint(id)})
	}
	return insertUser(t, db, u...)
}

func insertUser(t *testing.T, db *sql.DB, user ...user.User) func() {
	query := `INSERT INTO "users"
	(id, email, password, first_name, last_name, profile_image, created_at, updated_at) VALUES
	($1, $2, '', $3, $4, $5, $6, $7)`

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	for _, u := range user {
		_, err := stmt.Exec(
			u.ID,
			u.Email,
			u.FirstName,
			u.LastName,
			u.ProfileImage,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	return func() {
		defer db.Exec(`DELETE FROM users;`)
	}
}
