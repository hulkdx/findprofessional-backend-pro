package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func URLParamInt64(r *http.Request, key string) (int64, error) {
	str := chi.URLParam(r, key)
	return strconv.ParseInt(str, 10, 64)
}
