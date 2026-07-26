// Package service
package service

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/rizalta/url-short/server/internal/repo"
)

const (
	LENGTH     = 6
	CHARSET    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries = 5
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

	return "", errors.New("max retries reached for generating code")
}

func (s *service) GetURL(ctx context.Context, code string) (string, error) {
	u, err := s.queries.GetURL(ctx, code)
	if err != nil {
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
