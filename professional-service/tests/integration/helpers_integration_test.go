package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

func OutputSQL(t *testing.T, db *pgxpool.Pool, query string) {
	ctx := context.Background()

	rows, err := db.Query(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	// Collect column names
	fieldDescs := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescs))
	for i, fd := range fieldDescs {
		columns[i] = string(fd.Name)
	}
	fmt.Println(strings.Join(columns, "\t\t"))

	for rows.Next() {
		// rowValues is a slice of interface{}
		rowValues, err := rows.Values()
		if err != nil {
			t.Fatal(err)
		}

		var rowStrings []string
		for _, val := range rowValues {
			if val == nil {
				rowStrings = append(rowStrings, "NULL")
			} else {
				rowStrings = append(rowStrings, fmt.Sprintf("%v", val))
			}
		}
		fmt.Println(strings.Join(rowStrings, "\t\t"))
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func insertEmptyPro(t *testing.T, db *pgxpool.Pool) func() {
	return insertPro(t, db, professional.Professional{
		PriceNumber:   Int(0),
		PriceCurrency: String(""),
		Pending:       false,
	})
}

func insertPro(t *testing.T, db *pgxpool.Pool, pro ...professional.Professional) func() {
	ctx := context.Background()

	query := `INSERT INTO "professionals"
	(id,"email","password","first_name","last_name","coach_type","price_number","price_currency", "pending", "created_at", "updated_at") VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range pro {
		_, err := tx.Exec(ctx, query,
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
			tx.Rollback(ctx)
			t.Fatal(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	return func() {
		db.Exec(ctx, `DELETE FROM professionals`)
	}
}

func insertAvailability(t *testing.T, pool *pgxpool.Pool, availabilities ...professional.Availability) ([]int64, func()) {
	ctx := context.Background()

	query := `
        INSERT INTO "professional_availability"
            ("professional_id", "availability", "created_at", "updated_at")
        VALUES
            ($1, $2, $3, $4)
        RETURNING id
    `

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var ids []int64
	for _, a := range availabilities {
		availabilityRange := pgtype.Range[pgtype.Timestamptz]{
			Lower:     pgtype.Timestamptz{Time: a.From, Valid: true},
			Upper:     pgtype.Timestamptz{Time: a.To, Valid: true},
			LowerType: pgtype.Inclusive,
			UpperType: pgtype.Exclusive,
			Valid:     true,
		}

		var id int64
		execErr := tx.QueryRow(ctx, query,
			a.ProfessionalID,
			availabilityRange,
			a.CreatedAt,
			a.UpdatedAt,
		).Scan(&id)
		if execErr != nil {
			_ = tx.Rollback(ctx)
			t.Fatal(execErr)
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	return ids, func() {
		// Clean up both availability & professionals
		_, _ = pool.Exec(ctx, `DELETE FROM professional_availability; DELETE FROM professionals;`)
	}
}

// insertReview inserts Reviews in a single transaction.
func insertReview(t *testing.T, pool *pgxpool.Pool, review ...professional.Review) func() {
	ctx := context.Background()

	query := `
        INSERT INTO "professional_review"
            ("professional_id", "user_id", "rate", "created_at", "updated_at", "content_text", "id")
        VALUES
            ($1, $2, $3, $4, $5, $6, $7)
    `

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range review {
		_, execErr := tx.Exec(ctx, query,
			r.ProfessionalID,
			r.UserID,
			r.Rate,
			r.CreatedAt,
			r.UpdatedAt,
			r.ContentText,
			r.ID,
		)
		if execErr != nil {
			_ = tx.Rollback(ctx)
			t.Fatal(execErr)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	return func() {
		_, _ = pool.Exec(ctx, `DELETE FROM professional_review;`)
	}
}

// insertUserWithId inserts users that have their ID = email, for example.
func insertUserWithId(t *testing.T, pool *pgxpool.Pool, userIDs ...int) func() {
	users := []user.User{}
	for _, id := range userIDs {
		// example: Email set to string of the ID
		users = append(users, user.User{
			ID:    id,
			Email: fmt.Sprint(id),
		})
	}
	return insertUser(t, pool, users...)
}

// insertUser inserts one or more users in a single transaction.
func insertUser(t *testing.T, pool *pgxpool.Pool, users ...user.User) func() {
	ctx := context.Background()

	query := `
        INSERT INTO "users"
            (id, email, password, first_name, last_name, profile_image, created_at, updated_at)
        VALUES
            ($1, $2, '', $3, $4, $5, $6, $7)
    `

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range users {
		_, execErr := tx.Exec(ctx, query,
			u.ID,
			u.Email,
			u.FirstName,
			u.LastName,
			u.ProfileImage,
			time.Now(),
			time.Now(),
		)
		if execErr != nil {
			_ = tx.Rollback(ctx)
			t.Fatal(execErr)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	return func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users;`)
	}
}
