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

type service struct {
	queries *repo.Queries
}

func NewService(q *repo.Queries) *service {
	return &service{
		queries: q,
	}
}

func (s *service) Shorten(url string) (string, error) {
	for range maxRetries {
		code, err := generateCode()
		if err != nil {
			return "", err
		}

		res, err := s.queries.CreateURL(context.Background(), repo.CreateURLParams{Code: code, OriginalUrl: url})
		if err == nil {
			return res.Code, nil
		}
	}

	return "", errors.New("max retries reached for generating code")
}

func (s *service) GetURL(code string) string {
	u, err := s.queries.GetURL(context.Background(), code)
	if err != nil {
		return ""
	}

	return u
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
