package main

import (
	"firstgomode/internal/config"
	"firstgomode/internal/http-server/handlers/url/delete"
	"firstgomode/internal/http-server/handlers/url/redirect"
	"firstgomode/internal/http-server/handlers/url/save"
	mwLogger "firstgomode/internal/http-server/middleware/logger"
	"firstgomode/internal/lib/logger/handlers/slogpretty"
	"firstgomode/internal/lib/logger/sl"
	"firstgomode/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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
	log = log.With("env", cfg.Env)
	log.Info("starting project")
	log.Debug("debug logs are enabled")
	log.Error("error logs are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	//middleware

	router.Use(middleware.RequestID) // tracing
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/url/{alias}", delete.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Addres))

	srv := http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", sl.Err(err))
	}

	log.Error("server stopped")
	// TODO: запустить сервер
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog(env)
	case envDev:
		log = setupPrettySlog(env)
	case envProd:
		log = setupPrettySlog(env)
	default: // If env config is invalid, set prod settings by default due to security
		log = setupPrettySlog(env)
	}

	return log
}

func setupPrettySlog(env string) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	switch env {
	case envLocal:
		opts.SlogOpts.Level = slog.LevelDebug
	case envDev:
		opts.SlogOpts.Level = slog.LevelDebug
	case envProd:
		opts.SlogOpts.Level = slog.LevelInfo
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
