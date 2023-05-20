package utils

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var valid *validator.Validate
var lock sync.Mutex

func getValidator() *validator.Validate {
	lock.Lock()
	defer lock.Unlock()

	if valid == nil {
		valid = validator.New()
		return valid
	} else {
		return valid
	}
}

func IsValid(i interface{}) error {
	return getValidator().Struct(i)
}
