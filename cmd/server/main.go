package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalta/url-short/server/internal/handler"
	"github.com/rizalta/url-short/server/internal/repo"
	"github.com/rizalta/url-short/server/internal/service"
	"github.com/rizalta/url-short/server/web"
)

func main() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	schema := `
	CREATE TABLE IF NOT EXISTS urls (
 		code VARCHAR(10) PRIMARY KEY,
  	original_url TEXT NOT NULL,
  	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	if _, err := pool.Exec(ctx, schema); err != nil {
		log.Fatalf("faield to initialize database schema: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	staticMiddleware, err := web.StaticMiddleware()
	if err != nil {
		log.Fatalf("Failed to initialize static middleware: %v", err)
	}
	r.Use(staticMiddleware)

	q := repo.New(pool)
	s := service.NewService(q)
	h := handler.NewHandler(s)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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
