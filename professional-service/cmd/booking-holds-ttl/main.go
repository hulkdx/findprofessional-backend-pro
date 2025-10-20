package bookingholdsttl

import (
	"context"
	"os"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/db"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TIMEOUT = 2 * time.Minute
	LIMIT   = 1000
	GRACE   = 15 * time.Minute
)

func main() {
	logger.Debug("booking-holds-ttl started")

	timeout := TIMEOUT
	grace := GRACE
	limit := LIMIT

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	database := setupDatabase(ctx)
	defer db.Close(database)

	err := cleanup(ctx, database, &utils.RealTimeProvider{}, grace, limit)
	if err != nil {
		logger.Error("cleanup failed", err)
		os.Exit(1)
	}

	logger.Debug("booking-holds-ttl finished")
}

func setupDatabase(ctx context.Context) *pgxpool.Pool {
	databaseCfg := config.LoadDataBaseConfig()
	return db.Connect(ctx, databaseCfg)
}

func cleanup(
	ctx context.Context,
	db *pgxpool.Pool,
	timeProvider utils.TimeProvider,
	grace time.Duration,
	limit int,
) error {
	cutoff := timeProvider.Now().UTC().Add(-grace)
	for {
		query := `
			DELETE FROM booking_holds b
			USING (
				SELECT ctid FROM booking_holds
				WHERE expires_at < $1
				ORDER BY expires_at
				LIMIT $2
			) s
			WHERE b.ctid = s.ctid
		`
		res, err := db.Exec(ctx, query, cutoff, limit)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			break
		}
	}

	return nil
}
