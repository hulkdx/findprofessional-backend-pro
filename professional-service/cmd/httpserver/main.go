package main

import (
	"context"
	"errors"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	quit := make(chan os.Signal)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
		quit <- syscall.SIGTERM
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	// TODO: close db connection
	if err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server exiting")
}
