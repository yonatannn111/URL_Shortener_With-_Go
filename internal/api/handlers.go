package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type CreateReq struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
	ExpireDays int    `json:"expire_days,omitempty"`
}

type CreateResp struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (a *App) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// validate URL
	u := strings.TrimSpace(req.URL)
	parsed, err := url.ParseRequestURI(u)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.Error(w, "invalid url", http.StatusBadRequest); return
	}

	code := req.CustomCode
	if code == "" {
		code, _ = shortener.GenerateCode(6) // 6-char code
	}
	var expires *time.Time
	if req.ExpireDays > 0 {
		t := time.Now().Add(time.Duration(req.ExpireDays) * 24 * time.Hour)
		expires = &t
	}
	// Try inserting; if collision and generated, retry several times
	var id int
	for tries := 0; tries < 5; tries++ {
		id, err = a.Store.CreateURL(r.Context(), code, u, expires)
		if err != nil {
			// if unique violation -> generate new code and retry
			// check DB error text or error code (lib/pq) in real app
			code, _ = shortener.GenerateCode(6)
			continue
		}
		break
	}
	if err != nil {
		http.Error(w, "could not create short url", http.StatusInternalServerError); return
	}
	short := strings.TrimRight(a.BaseURL, "/") + "/" + code
	json.NewEncoder(w).Encode(CreateResp{Code: code, ShortURL: short})
}

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	ctx := r.Context()

	rec, err := a.Store.GetURLByCode(ctx, code)
	if err != nil || rec == nil {
		http.NotFound(w, r); return
	}
	// Check expiry if present
	if rec.ExpiresAt != nil && rec.ExpiresAt.Before(time.Now()) {
		http.NotFound(w, r); return
	}

	// Extract IP (handle proxies)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	ua := r.UserAgent()
	ref := r.Referer()

	// async: geo lookup (best-effort)
	country, city := lookupGeo(ip) // example function below

	// push event to click worker
	a.Worker.ClickCh <- worker.ClickEvent{
		URLID: rec.ID, Code: rec.Code, IP: ip, Country: country, City: city, UserAgent: ua, Referer: ref,
	}

	http.Redirect(w, r, rec.Original, http.StatusTemporaryRedirect)
}

func (a *App) StatsHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	rec, err := a.Store.GetURLByCode(r.Context(), code)
	if err != nil || rec == nil {
		http.NotFound(w, r); return
	}
	json.NewEncoder(w).Encode(rec)
}
