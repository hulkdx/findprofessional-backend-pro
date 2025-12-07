package professionalrepo

import (
	"context"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
)

func (r *RepositoryImpl) GetBookingStatus(ctx context.Context, bookingId int64) (professional.StatusResponse, error) {
	const query = `SELECT status FROM bookings WHERE id = $1`
	var result professional.StatusResponse
	err := r.db.QueryRow(ctx, query, bookingId).Scan(&result.Status)
	if err != nil {
		return professional.StatusResponse{}, err
	}
	return result, nil
}
