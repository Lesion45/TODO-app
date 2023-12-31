package main

import (
	"TODOapp/internal/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("address", cfg.Server.Address))
	log.Debug("logger debug mode enabled")

	//TODO: STORAGE

	//TODO: ROUTER
	router := chi.NewRouter()

	//TODO: RUN SERVER
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.TimeOut,
		WriteTimeout: cfg.Server.TimeOut,
		IdleTimeout:  cfg.Server.IdleTimeOut,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
