package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToTstzrange(from time.Time, to time.Time) *pgtype.Range[pgtype.Timestamp] {
	return &pgtype.Range[pgtype.Timestamp]{
		Lower: pgtype.Timestamp{
			Time:             from,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		Upper: pgtype.Timestamp{
			Time:             to,
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}
}
