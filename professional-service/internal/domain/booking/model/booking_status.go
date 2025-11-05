package bookingmodel

type BookingStatus string

// From this link: https://sameerahmed56.medium.com/an-uncomfortably-deep-dive-into-the-idempotency-key-67626c8d3f3d
const (
	BookingStatusNonExistent BookingStatus = "non-existent"
	BookingStatusProcessing  BookingStatus = "processing"
	BookingStatusCompleted   BookingStatus = "completed"
	BookingStatusFailed      BookingStatus = "failed"
	BookingStatusExpired     BookingStatus = "expired"
)
