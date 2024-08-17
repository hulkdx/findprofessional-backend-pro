package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/user"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func CreateProTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))

	t.Run("not found a record, create a new record", func(t *testing.T) {
		// Arrange
		bodyRequest := professional.CreateRequest{
			Email:         "test@gmail.com",
			Password:      "P@ssw0rd123",
			FirstName:     "John",
			LastName:      "Doe",
			SkypeId:       "john_doe_skype",
			AboutMe:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Price:         500,
			PriceCurrency: "USD",
			CoachType:     "Lifecoach",
		}
		defer db.Exec(`DELETE FROM professionals`)

		user := user.User{
			ID:        1,
			Email:     bodyRequest.Email,
			FirstName: bodyRequest.FirstName,
			LastName:  bodyRequest.LastName,
		}
		d0 := insertUser(t, db, user)
		defer d0()

		request := createProRequest(bodyRequest)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusCreated)
	})

	t.Run("not found a record, pending is stored as true in the database", func(t *testing.T) {
		// Arrange
		bodyRequest := professional.CreateRequest{
			Email:         "test@gmail.com",
			Password:      "P@ssw0rd123",
			FirstName:     "John",
			LastName:      "Doe",
			SkypeId:       "john_doe_skype",
			AboutMe:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Price:         500,
			PriceCurrency: "USD",
			CoachType:     "Lifecoach",
		}
		defer db.Exec(`DELETE FROM professionals`)

		user := user.User{
			ID:        1,
			Email:     bodyRequest.Email,
			FirstName: bodyRequest.FirstName,
			LastName:  bodyRequest.LastName,
		}
		d0 := insertUser(t, db, user)
		defer d0()

		request := createProRequest(bodyRequest)
		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, request)
		// Asserts
		assert.Equal(t, response.Code, http.StatusCreated)
		result := getPendingFromDatabase(db)
		assert.Equal(t, *result, true)
	})

	t.Run("found a record, then return 409 status", func(t *testing.T) {
		// Arrange
		bodyRequest := professional.CreateRequest{
			Email:         "test@gmail.com",
			Password:      "P@ssw0rd123",
			FirstName:     "John",
			LastName:      "Doe",
			SkypeId:       "john_doe_skype",
			AboutMe:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Price:         500,
			PriceCurrency: "USD",
			CoachType:     "Lifecoach",
		}
		defer db.Exec(`DELETE FROM professionals`)

		user := user.User{
			ID:        1,
			Email:     bodyRequest.Email,
			FirstName: bodyRequest.FirstName,
			LastName:  bodyRequest.LastName,
		}
		d0 := insertUser(t, db, user)
		defer d0()

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, createProRequest(bodyRequest))
		assert.Equal(t, response.Code, http.StatusCreated)

		response = httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, createProRequest(bodyRequest))
		// Asserts
		assert.Equal(t, response.Code, http.StatusConflict)
	})

	t.Run("Creating should update user.professionalId", func(t *testing.T) {
		// Arrange
		bodyRequest := professional.CreateRequest{
			Email:         "test@gmail.com",
			Password:      "P@ssw0rd123",
			FirstName:     "John",
			LastName:      "Doe",
			SkypeId:       "john_doe_skype",
			AboutMe:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			Price:         500,
			PriceCurrency: "USD",
			CoachType:     "Lifecoach",
		}
		defer db.Exec(`DELETE FROM professionals`)

		user := user.User{
			ID:        1,
			Email:     bodyRequest.Email,
			FirstName: bodyRequest.FirstName,
			LastName:  bodyRequest.LastName,
		}
		d0 := insertUser(t, db, user)
		defer d0()

		response := httptest.NewRecorder()
		// Act
		handler.ServeHTTP(response, createProRequest(bodyRequest))
		// Asserts
		assert.Equal(t, response.Code, http.StatusCreated)

		var responseProfessionalId int64 = -1
		db.QueryRow("SELECT professional_id from users WHERE email=$1", bodyRequest.Email).Scan(&responseProfessionalId)
		assert.NotEqual(t, responseProfessionalId, -1)
	})
}

func createProRequest(body professional.CreateRequest) *http.Request {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	return NewJsonRequest("PUT", "/professional", &buf)
}

func getPendingFromDatabase(db *sql.DB) *bool {
	var pending sql.NullBool
	db.QueryRow("SELECT pending FROM professionals").Scan(&pending)
	return &pending.Bool
}
