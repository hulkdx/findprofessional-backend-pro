package utils

import (
	"errors"
	"net/http"
)

var ErrUnknown = errors.New("unknown")

var ErrDuplicate = errors.New("duplicate")

var ErrNotFoundUser = errors.New("user not found")

var ErrIdempotencyKeyIsUsed = errors.New("idempotency key is used")

var ErrAvailabilityOwnershipMismatch = errors.New("id does not belong to professional")

var ErrAvailabilityDoesNotExist = NewHttpError("availability_id doesn't exist", http.StatusNotFound)
