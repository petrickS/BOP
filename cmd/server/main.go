package main

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/petrickS/BOP/internal/core/config"
	"github.com/petrickS/BOP/internal/core/db"
	"github.com/petrickS/BOP/internal/features/search"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	s := &server{
		provider:  search.NewPostgresProvider(pool),
		templates: tmpl,
	}

	r := chi.NewRouter()
	r.Get("/", s.handleIndex)
	r.Post("/search", s.handleSearchHTML)
	r.Post("/api/search", s.handleSearchJSON)

	httpServer := &http.Server{Addr: ":8080", Handler: r}

	serveErr := make(chan error, 1)
	go func() {
		log.Println("listening on :8080")
		serveErr <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutting down...")
		stop() // перестаём перехватывать сигналы, повторный Ctrl+C завершит процесс сразу

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}

	log.Println("stopped")
}
