package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/db"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

func main() {
	logger.Debug("Server started")

	ctx := context.Background()
	cfg := config.Load()

	database := db.Connect(ctx, cfg.Database)
	defer db.Close(database)

	handler := router.NewHandler(database)
	server := newServer(ctx, cfg, handler)
	<-listenAndServe(server)
	err := shutdown(server)
	if err != nil {
		logger.Error("Server stopped error:", err)
	} else {
		logger.Debug("Server stopped")
	}
}

func newServer(ctx context.Context, cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		BaseContext:  func(net.Listener) context.Context { return ctx },
	}
}

func listenAndServe(server *http.Server) <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
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
