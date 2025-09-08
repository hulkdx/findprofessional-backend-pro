package utils

import "errors"

var ErrUnknown = errors.New("unknown")

var ErrDuplicate = errors.New("duplicate")

var ErrNotFoundUser = errors.New("user not found")

var ErrAmountInCentsMismatch = errors.New("amount_in_cents mismatch")

var ErrCurrencyMismatch = errors.New("currency mismatch")

var ErrValidationDatabase = errors.New("validation error: 101")
