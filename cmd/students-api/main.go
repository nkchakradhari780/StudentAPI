package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/nkchakradhari780/students-api/internal/config"
	"github.com/nkchakradhari780/students-api/internal/http/handlers/student"
)

func main() {
	//load config
	cfg := config.MustLoad()
	// database setup
	// router setup
	router := http.NewServeMux()
	
	router.HandleFunc("POST /api/students", student.New())
	// server setup

	server := http.Server{
		Addr: cfg.HttpServer.Addr,
		Handler: router,
	}

	slog.Info("Starting the server", slog.String("address", cfg.HttpServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to Start server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Stopped")
}
