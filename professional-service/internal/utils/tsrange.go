package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertToTsRange(date string, from string, to string) (*pgtype.Range[pgtype.Timestamp], error) {
	lowerTime, err := time.Parse("2006-01-02 15:00", fmt.Sprintf("%s %s", date, from))
	if err != nil {
		return nil, err
	}
	upperTime, err := time.Parse("2006-01-02 15:00", fmt.Sprintf("%s %s", date, to))
	if err != nil {
		return nil, err
	}

	return &pgtype.Range[pgtype.Timestamp]{
		Lower: pgtype.Timestamp{
			Time:             lowerTime,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		Upper: pgtype.Timestamp{
			Time:             upperTime,
			Valid:            true,
			InfinityModifier: pgtype.Finite,
		},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
		Valid:     true,
	}, nil
}
