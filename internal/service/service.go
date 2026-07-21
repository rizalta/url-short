// Package service
package service

import (
	"crypto/rand"
	"errors"
	"sync"
)

const (
	LENGTH     = 6
	CHARSET    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries = 5
)

type service struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewService() *service {
	return &service{
		urls: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (s *service) Shorten(url string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for range maxRetries {
		code, err := generateCode()
		if err != nil {
			return "", err
		}
		if _, ok := s.urls[code]; !ok {
			s.urls[code] = url
			return code, nil
		}
	}

	return "", errors.New("max retries reached for generating code")
}

func (s *service) GetURL(code string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.urls[code]
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
