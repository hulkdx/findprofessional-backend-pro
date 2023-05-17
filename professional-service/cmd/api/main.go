package main

import (
	"context"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/domain/professional"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/db"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func main() {
	logger.Debug("Server started")
	cfg := config.Load()

	database := db.Connect(cfg.Database)
	defer db.Close(database)

	handler := router.Handler(professional.NewControllerFromDB(database))
	server := newServer(cfg, handler)
	<-listenAndServe(server)
	err := shutdown(server)
	if err != nil {
		logger.Error("Server stopped error:", err)
	} else {
		logger.Debug("Server stopped")
	}
}

func newServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
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
