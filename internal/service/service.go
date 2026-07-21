// Package service
package service

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
)

const (
	LENGTH     = 6
	CHARSET    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries = 5
)

type service struct {
	urls map[string]string
	mu   *sync.Mutex
}

func NewService() *service {
	return &service{
		urls: make(map[string]string),
		mu:   &sync.Mutex{},
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
	return s.urls[code]
}

func generateCode() (string, error) {
	b := make([]byte, LENGTH)
	charsetLen := big.NewInt(int64(len(CHARSET)))

	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = CHARSET[n.Int64()]
	}

	return string(b), nil
}
