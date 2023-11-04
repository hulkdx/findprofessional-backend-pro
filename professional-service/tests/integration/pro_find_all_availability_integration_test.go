package integration_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/civil"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/tests/assert"
)

func FindAllAvailabilityProfessionalTest(t *testing.T, db *sql.DB) {
	handler := router.Handler(NewTestController(db))
	t.Run("empty availability", func(t *testing.T) {
		// Arrange
		expected := []professional.Availability{}
		_, err := db.Exec(`INSERT INTO "professionals"
		(id,"email","password","first_name","last_name","coach_type","price_number","price_currency") VALUES
		(1 ,''			,''				 ,''          ,''          ,''          ,0					  ,''				      )
		`)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Exec(`DELETE FROM professional_availability; DELETE FROM professionals;`)

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)

		assert.Equal(t, len(response_model), 1)
		assert.EqualAnyOrder(t, response_model[0].Availability, expected)
	})
	t.Run("testing adding some availabilities", func(t *testing.T) {
		// Arrange
		expected := []professional.Availability{
			{
				ID:   1,
				Date: civil.Date{Year: 2023, Month: 11, Day: 4},
				From: civil.Time{Hour: 5, Minute: 30},
				To:   civil.Time{Hour: 6, Minute: 30},
			},
			{
				ID:   2,
				Date: civil.Date{Year: 2020, Month: 11, Day: 4},
				From: civil.Time{Hour: 15, Minute: 30},
				To:   civil.Time{Hour: 16, Minute: 00},
			},
		}

		_, err := db.Exec(`INSERT INTO "professionals"
		(id,"email","password","first_name","last_name","coach_type","price_number","price_currency") VALUES
		(1 ,''			,''				 ,''          ,''          ,''          ,0					  ,''				      )
		`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(`INSERT INTO "professional_availability"
		("professional_id","date"				 ,"from"    ,"to"      ) VALUES
		(1								,'"2023-11-04"','05:30:00','06:30:00'),
		(1								,'"2020-11-04"','15:30:00','16:00:00')
		`)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Exec(`DELETE FROM professional_availability; DELETE FROM professionals;`)

		request := NewJsonRequest("GET", "/professional", nil)
		response := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(response, request)

		// Assert
		assert.Equal(t, response.Code, http.StatusOK)

		response_model := []professional.Professional{}
		Unmarshal(response, &response_model)

		assert.Equal(t, len(response_model), 1)
		assert.EqualAnyOrder(t, response_model[0].Availability, expected)
	})
}
