package booking_model

type BookingStatus string

const (
	BookingStatusHold      BookingStatus = "hold"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCanceled  BookingStatus = "canceled"
)
