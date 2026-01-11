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
	(
		"id",
		"email",
		"password",
		"first_name",
		"last_name",
		"coach_type",
		"price_number",
		"price_currency",
		"pending",
		"session_platform",
		"session_link",
		"created_at",
		"updated_at"
	) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

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
			p.SessionPlatform,
			p.SessionLink,
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
func insertUserWithId(t *testing.T, pool *pgxpool.Pool, userIDs ...int64) func() {
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

type TestBookingItems struct {
	BookingID      int64
	AvailabilityID int64
}

func insertBooking(
	t *testing.T,
	pool *pgxpool.Pool,
	userId,
	proId int64,
	status,
	currency,
	paymentIntent string,
	scheduledStartAt *time.Time,
	scheduledEndAt *time.Time,
) (int64, func()) {
	return insertBookingWithSessions(
		t,
		pool,
		userId,
		proId,
		status,
		currency,
		paymentIntent,
		scheduledStartAt,
		scheduledEndAt,
		nil,
		nil,
	)
}

func insertBookingWithSessions(
	t *testing.T,
	pool *pgxpool.Pool,
	userId,
	proId int64,
	status,
	currency,
	paymentIntent string,
	scheduledStartAt *time.Time,
	scheduledEndAt *time.Time,
	sessionLink *string,
	sessionPlatform *string,
) (int64, func()) {
	t.Helper()
	ctx := context.Background()

	query := `INSERT INTO bookings (
			user_id,
			professional_id,
			status,
			total_amount_cents,
			currency,
			stripe_payment_intent_id,
			scheduled_start_at,
			scheduled_end_at,
			session_platform,
			session_link,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now(), now())
		RETURNING id`

	// use some fixed amount for tests, or adjust if you want it parametrized
	const totalAmountCents int64 = 1000

	var id int64
	if err := pool.QueryRow(ctx, query,
		userId,
		proId,
		status,
		totalAmountCents,
		currency,
		paymentIntent,
		scheduledStartAt,
		scheduledEndAt,
		sessionPlatform,
		sessionLink,
	).Scan(&id); err != nil {
		t.Fatalf("failed to insert booking: %v", err)
	}

	return id, func() {
		pool.Exec(ctx, `DELETE FROM bookings`)
	}
}

func insertBookingItems(t *testing.T, pool *pgxpool.Pool, items ...TestBookingItems) func() {
	ctx := context.Background()

	query := `INSERT INTO booking_items (booking_id, availability_id, created_at, updated_at) VALUES ($1, $2, now(), now())`

	for _, it := range items {
		_, err := pool.Exec(ctx, query,
			it.BookingID,
			it.AvailabilityID,
		)
		if err != nil {
			t.Fatalf("failed to insert booking_item: %+v, err: %v", it, err)
		}
	}

	return func() {
		_, _ = pool.Exec(ctx, `DELETE FROM booking_items;`)
	}
}
