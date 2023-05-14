package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/db"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/logger"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
)

func main() {
	logger.Debug("Server started")
	cfg := config.Load()

	database := db.Connect(cfg.Database)
	defer database.Close()

	server := newServer(cfg)
	<-listenAndServe(server)
	err := shutdown(server)
	if err != nil {
		logger.Error("Server stopped error:", err)
	} else {
		logger.Debug("Server stopped")
	}
}

func newServer(cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func listenAndServe(server *http.Server) chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Error("ListenAndServe", err)
		}
		quit <- syscall.SIGTERM
	}()
	return quit
}

func shutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}
