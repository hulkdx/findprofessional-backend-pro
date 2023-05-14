package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/config"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/db"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/router"
)

func main() {
	log.Println("Server starting")
	cfg := config.Load()

	database := db.Connect(cfg.Database)
	defer database.Close()

	server := Server(cfg)
	<-ListenAndServe(server)
	Shutdown(server)
}

func Server(cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      router.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func ListenAndServe(server *http.Server) chan os.Signal {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Server error: %v\n", err)
		}
		quit <- syscall.SIGTERM
	}()
	return quit
}

func Shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Server shutdown with error:%v\n", err)
	} else {
		log.Printf("Server shutdown")
	}
}
