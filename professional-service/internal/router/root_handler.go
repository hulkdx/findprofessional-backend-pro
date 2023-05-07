package router

import (
	"net/http"
)

func Handler() http.Handler {
	return ChiRouter()
}
