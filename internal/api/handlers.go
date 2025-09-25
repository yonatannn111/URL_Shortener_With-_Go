package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/worker"
)

// App holds store, worker, and base URL
type App struct {
	Store   *storage.Store
	Worker  *worker.Worker
	BaseURL string
}

// generateCode generates a random n-character string
func generateCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ShortenHandler creates a short link and records analytics
func (a *App) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL cannot be empty"})
		return
	}

	// Generate unique code
	code := generateCode(6)

	// Save URL in store
	if err := a.Store.SaveURL(code, req.URL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to save URL"})
		return
	}

	// Record analytics asynchronously
	go func(ip string) {
		ip = strings.Split(ip, ":")[0] // remove port if present
		if ip == "127.0.0.1" || ip == "::1" {
			return // skip local IPs
		}
		country, city := lookupGeo(ip)
		a.Worker.AddAnalytics(code, ip, country, city)
	}(r.RemoteAddr)

	resp := map[string]string{"short_url": a.BaseURL + "/" + code}
	json.NewEncoder(w).Encode(resp)
}

// RedirectHandler redirects short link and records click
func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, err := a.Store.GetURL(code)

	if err != nil || url == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "URL not found"})
		return
	}

	// Record click asynchronously
	go func(ip string) {
		ip = strings.Split(ip, ":")[0] // remove port
		if ip == "127.0.0.1" || ip == "::1" {
			return // skip local IPs
		}
		country, city := lookupGeo(ip)
		a.Worker.AddAnalytics(code, ip, country, city)
	}(r.RemoteAddr)

	http.Redirect(w, r, url, http.StatusFound)
}

// StatsHandler returns total clicks for a code
func (a *App) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code := chi.URLParam(r, "code")
	clicks := a.Worker.GetClicks(code)

	resp := map[string]interface{}{
		"code":   code,
		"clicks": clicks,
	}

	json.NewEncoder(w).Encode(resp)
}
