package main

import (
	"fmt"
	"github.com/hulkdx/findprofessional-backend-pro/api"
	"net/http"
	"time"
)

func main() {
	// TODO: read db connection env
	// TODO: database connection
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      &api.Router{},
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	server.ListenAndServe()
	// TODO: shutdown gracefully
}
