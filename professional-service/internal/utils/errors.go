package utils

import "errors"

var ErrUnknown = errors.New("unknown")

var ErrDuplicate = errors.New("duplicate")

var ErrNotFoundUser = errors.New("user not found")

var ErrIdempotencyKeyIsUsed = errors.New("idempotency key is used")

var ErrAvailabilityOwnershipMismatch = errors.New("id does not belong to professional")
