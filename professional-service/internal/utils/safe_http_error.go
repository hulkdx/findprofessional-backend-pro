package utils

type SafeHttpError struct {
	msg string
}

func NewSafeHttpError(msg string) *SafeHttpError {
	return &SafeHttpError{msg: msg}
}

func (e *SafeHttpError) Error() string { return e.msg }
