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
	"github.com/nkchakradhari780/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()
	// database setup

	storage , err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version","1.0.0"))

	// router setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.CreateNewStudent(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
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
