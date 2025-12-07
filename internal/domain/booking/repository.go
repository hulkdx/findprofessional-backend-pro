package booking

import "context"

type Repository interface {
	GetStatus(ctx context.Context, bookingId string) (string, error)
}
