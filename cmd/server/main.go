package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/petrickS/BOP/internal/core/config"
	"github.com/petrickS/BOP/internal/core/db"
	"github.com/petrickS/BOP/internal/features/search"
)

//go:embed templates/*.html
var templatesFS embed.FS

type server struct {
	provider  search.Provider
	templates *template.Template
}

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

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("render index: %v", err)
	}
}

func (s *server) handleSearchHTML(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	qty, err := strconv.Atoi(r.FormValue("qty"))
	if err != nil {
		http.Error(w, "invalid qty", http.StatusBadRequest)
		return
	}

	req := search.Request{
		Product: r.FormValue("product"),
		Qty:     qty,
		Region:  r.FormValue("region"),
	}

	result, err := search.Handle(r.Context(), s.provider, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.templates.ExecuteTemplate(w, "results.html", result); err != nil {
		log.Printf("render results: %v", err)
	}
}

func (s *server) handleSearchJSON(w http.ResponseWriter, r *http.Request) {
	var req search.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := search.Handle(r.Context(), s.provider, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("encode response: %v", err)
	}
}
