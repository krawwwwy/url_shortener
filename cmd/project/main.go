package main

import (
	"firstgomode/internal/config"
	"firstgomode/internal/lib/config/sl"
	"firstgomode/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// fmt.Println(cfg) //после дебага удалить в проде в конфиге серьезные вещи

	log := setupLogger(cfg.Env)

	log.Info("starting project", "env", cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()
	_ = router

	//middleware

	//TODO: Инициализировать (подключить) router

	// TODO: запустить сервер
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
