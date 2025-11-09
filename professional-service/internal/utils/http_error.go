package utils

type HttpError struct {
	msg    string
	Status int
}

func NewHttpError(msg string, status int) *HttpError {
	return &HttpError{
		msg:    msg,
		Status: status,
	}
}

func (e *HttpError) Error() string { return e.msg }
