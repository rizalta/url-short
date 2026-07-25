package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rizalta/url-short/server/internal/handler"
	"github.com/rizalta/url-short/server/internal/service"
	"github.com/rizalta/url-short/server/web"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	staticMiddleware, err := web.StaticMiddleware()
	if err != nil {
		log.Fatalf("Failed to initialize static middleware: %v", err)
	}
	r.Use(staticMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := service.NewService()
	h := handler.NewHandler(s)
	r.Mount("/", h.Routes())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
