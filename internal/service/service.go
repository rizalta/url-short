// Package service
package service

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rizalta/url-short/server/internal/repo"
)

const (
	LENGTH     = 6
	CHARSET    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries = 5
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrCodeFailed  = errors.New("code generation failed")
)

type Querier interface {
	CreateURL(ctx context.Context, params repo.CreateURLParams) (repo.Url, error)
	GetURL(ctx context.Context, code string) (string, error)
}

type service struct {
	queries Querier
}

func NewService(q Querier) *service {
	return &service{
		queries: q,
	}
}

func (s *service) Shorten(ctx context.Context, url string) (string, error) {
	for range maxRetries {
		code, err := generateCode()
		if err != nil {
			return "", err
		}

		res, err := s.queries.CreateURL(ctx, repo.CreateURLParams{Code: code, OriginalUrl: url})
		if err == nil {
			return res.Code, nil
		}
	}

	return "", ErrCodeFailed
}

func (s *service) GetURL(ctx context.Context, code string) (string, error) {
	u, err := s.queries.GetURL(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	return u, nil
}

func generateCode() (string, error) {
	b := make([]byte, LENGTH)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = CHARSET[b[i]%byte(len(CHARSET))]
	}

	return string(b), nil
}
