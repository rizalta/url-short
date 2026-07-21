package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rizalta/url-short/server/internal/handler"
	"github.com/rizalta/url-short/server/internal/service"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := service.NewService()
	h := handler.NewHandler(s)

	r.Mount("/", h.Routes())

	_ = http.ListenAndServe(":8000", r)
}
