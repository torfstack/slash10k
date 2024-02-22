package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"scurvy10k/src/handler"
	"time"
)

func main() {
	setupLogger()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	slog.Debug("setting up routes")

	r.Route("/api", func(r chi.Router) {
		r.Get("/debt", handler.Debt)

		r.Route("/chars", func(r chi.Router) {
			r.Post("/", handler.AddChar)
			r.Get("/{id}", handler.GetChar)
			r.Delete("/{id}", handler.DeleteChar)
		})
	})

	slog.Info("server started", "port", 3000)
	slog.Error("server stopped", http.ListenAndServe(":3000", r))
}

func setupLogger() {
	level := new(slog.LevelVar)
	level.Set(slog.LevelDebug) // change this to change log level
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
