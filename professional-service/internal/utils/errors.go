package utils

import "errors"

var ErrUnknown = errors.New("unknown")

var ErrDuplicate = errors.New("duplicate")

var ErrNotFoundUser = errors.New("user not found")

var ErrIdempotencyKeyExpired = errors.New("idempotency key expired")
