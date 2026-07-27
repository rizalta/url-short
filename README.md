# URL Shortener

A full-stack URL shortener application built to learn modern cloud deployment practices on an Oracle Cloud Infrastructure (OCI) server.

Live Application: https://url-shortz.duckdns.org

## What Was Implemented

This project was built incrementally to study containerization, networking, databases, CI/CD pipelines, and production hardening.

- Single Binary Architecture: Integrated a React, Vite, and TypeScript frontend into the Go server using `go:embed`. The server serves both the UI and API from a single compiled binary without CORS issues.
- Reverse Proxy and HTTPS: Configured Caddy on OCI to handle incoming traffic on ports 80 and 443, automatically provisioning and renewing Let's Encrypt SSL/TLS certificates.
- Database Persistence: Integrated PostgreSQL 16 managed via Docker Compose. Database interactions use `sqlc` to generate type-safe Go code from raw SQL queries.
- Automated CI/CD: Built a GitHub Actions pipeline that runs backend tests, builds the frontend, and deploys updates to the OCI server over SSH on every push to `main`.
- Error Handling and Observability: Implemented domain-specific sentinel errors, structured JSON error responses, database-aware `/health` checks, and container log rotation.

## Tech Stack

- Backend: Go 1.26, Chi, sqlc, pgx/v5
- Frontend: React, Vite, TypeScript, Catppuccin Mocha theme
- Database: PostgreSQL 16
- Infrastructure: Docker Compose, Caddy, OCI Free Tier
- CI/CD: GitHub Actions

## Local Development

1. Start PostgreSQL:
   ```bash
   docker compose up -d postgres
   ```

2. Build the frontend:
   ```bash
   cd web && pnpm install && pnpm build && cd ..
   ```

3. Run the Go server:
   ```bash
   go run ./cmd/server/main.go
   ```

## Running Tests

Run unit tests across backend services and handlers:

```bash
go test -v ./...
```
