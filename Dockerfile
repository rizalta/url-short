# Frontend
FROM node:23-alpine AS frontend-builder
WORKDIR /app/web

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY web/ ./
RUN pnpm build

# Backend
FROM golang:1.26.5-alpine AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime
FROM alpine:latest
COPY --from=backend-builder /app/server /server

EXPOSE 8000

ENTRYPOINT ["/server"]
