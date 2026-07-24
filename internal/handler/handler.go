// Package handler
package handler

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

var ErrInvalidURL = errors.New("invalid URL")

type Service interface {
	Shorten(string) (string, error)
	GetURL(string) string
}

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{s}
}

type ShortenReq struct {
	URL string `json:"url"`
}

type ShortenRes struct {
	ShortURL string `json:"short_url"`
}

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	var req ShortenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsed, err := normalizeURL(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := h.service.Shorten(parsed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ShortenRes{code}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if code == "" {
		http.Error(w, "code required", http.StatusBadRequest)
		return
	}

	u := h.service.GetURL(code)
	if u == "" {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, u, http.StatusFound)
}

func (h *handler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/shorten", h.shorten)
	r.Get("/{code}", h.getURL)
	return r
}

func normalizeURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return "", ErrInvalidURL
	}

	if !strings.HasPrefix(rawURL, "https://") && !strings.HasPrefix(rawURL, "http://") {
		rawURL = "https://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	host := u.Hostname()
	if host == "" {
		return "", ErrInvalidURL
	}

	if !isValidHost(host) {
		return "", ErrInvalidURL
	}

	return rawURL, nil
}

func isValidHost(host string) bool {
	if host == "localhost" || net.ParseIP(host) != nil {
		return true
	}

	if !strings.Contains(host, ".") || strings.HasPrefix(host, ".") || strings.HasSuffix(host, ".") {
		return false
	}

	return true
}
