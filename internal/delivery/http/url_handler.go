package http

import (
	"encoding/json"
	"html/template"
	"net/http"

	"url-shortener/internal/domain"

	"github.com/go-chi/chi/v5"
)

type urlHandler struct {
	usecase domain.URLUsecase
}

// NewURLHandler mounts all URL routes
func NewURLHandler(r *chi.Mux, u domain.URLUsecase) {
	handler := &urlHandler{
		usecase: u,
	}

	// API Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.HealthCheck)
		r.Post("/urls", handler.ShortenURL)
		r.Get("/urls/{shortCode}/stats", handler.GetStats)
	})

	// Redirect Route
	r.Get("/{shortCode}", handler.Redirect)
	
	// Frontend SSR Route
	r.Get("/", handler.RenderIndex)
}

func (h *urlHandler) RenderIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (h *urlHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
}

func (h *urlHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OriginalURL string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
		return
	}

	urlObj, err := h.usecase.ShortenURL(r.Context(), req.OriginalURL)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(urlObj)
}

func (h *urlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	originalURL, err := h.usecase.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (h *urlHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	stats, err := h.usecase.GetURLStats(r.Context(), shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
