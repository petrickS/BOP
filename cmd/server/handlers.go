package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/petrickS/BOP/internal/features/search"
)

//go:embed templates/*.html
var templatesFS embed.FS

type server struct {
	provider  search.Provider
	templates *template.Template
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
