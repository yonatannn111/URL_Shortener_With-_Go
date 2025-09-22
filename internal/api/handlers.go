package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/worker"
)

type App struct {
	Store   *storage.Store
	Worker  *worker.Worker
	BaseURL string
}

// ShortenHandler - create short link
func (a *App) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	code := "abc123" // TODO: generate unique code
	if err := a.Store.SaveURL(r.Context(), code, req.URL); err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"short_url": a.BaseURL + "/" + code}
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler - redirect short link
func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, err := a.Store.GetURL(r.Context(), code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// StatsHandler - dummy stats for now
func (a *App) StatsHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	resp := map[string]string{
		"code":  code,
		"click": "42", // placeholder
	}
	json.NewEncoder(w).Encode(resp)
}
