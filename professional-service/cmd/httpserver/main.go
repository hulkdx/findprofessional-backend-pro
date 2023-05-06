package main

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"net/http"
)

func main() {
	cfg := config.Load()
	// TODO: read db connection env
	// TODO: database connection
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      &internal.Router{},
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	server.ListenAndServe()
	// TODO: shutdown gracefully
}
