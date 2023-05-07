package main

import (
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"net/http"
)

func main() {
	cfg := config.Load()
	// TODO: database connection
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	server.ListenAndServe()
	// TODO: shutdown gracefully
}
