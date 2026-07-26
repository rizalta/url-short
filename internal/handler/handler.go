// Package handler
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rizalta/url-short/server/internal/service"
)

var ErrInvalidURL = errors.New("invalid URL")

type Service interface {
	Shorten(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, error)
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
	Code string `json:"code"`
}

type ErrorRes struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorRes{Error: message}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()
	var req ShortenReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	parsed, err := normalizeURL(req.URL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	code, err := h.service.Shorten(r.Context(), parsed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	res := ShortenRes{Code: code}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func (h *handler) getURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if code == "" {
		respondWithError(w, http.StatusBadRequest, "Code Required")
		return
	}

	u, err := h.service.GetURL(r.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			respondWithError(w, http.StatusNotFound, "URL not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		}
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
