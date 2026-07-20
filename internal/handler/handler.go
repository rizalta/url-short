// Package handler
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Service interface{}

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{s}
}

func (h *handler) shorten(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/shorten", h.shorten)
	return r
}
