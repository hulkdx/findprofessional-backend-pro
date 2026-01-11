package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToTstzrange(from time.Time, to time.Time) *pgtype.Range[pgtype.Timestamptz] {
	return &pgtype.Range[pgtype.Timestamptz]{
		Lower: pgtype.Timestamptz{
			Time:             from,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		Upper: pgtype.Timestamptz{
			Time:             to,
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}
}
