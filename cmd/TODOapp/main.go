package main

import (
	"TODOapp/internal/config"
	"TODOapp/internal/http-server/handlers/task/add"
	st "TODOapp/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("address", cfg.Server.Address))
	log.Debug("logger debug mode enabled")

	log.Info("initializing redis")
	redisClient := st.NewRedisClient()
	err := redisClient.IsRunning()
	if err != nil {
		log.Error("failed to initialize storage", err)
		os.Exit(1)
	}
	log.Info("connected to redis")

	//TODO: ROUTERS
	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Post("/task", add.New(log, redisClient))

	log.Info("server is starting")
	server := &http.Server{
		Addr:         "localhost:2000",
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
