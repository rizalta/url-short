package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rizalta/url-short/server/internal/repo"
)

type mockQuerier struct {
	urls map[string]string
}

func newMockQuerier() *mockQuerier {
	return &mockQuerier{urls: make(map[string]string)}
}

func (m *mockQuerier) CreateURL(ctx context.Context, params repo.CreateURLParams) (repo.Url, error) {
	if _, ok := m.urls[params.Code]; ok {
		return repo.Url{}, errors.New("code collision")
	}
	m.urls[params.Code] = params.OriginalUrl
	return repo.Url{Code: params.Code, OriginalUrl: params.OriginalUrl}, nil
}

func (m *mockQuerier) GetURL(ctx context.Context, code string) (string, error) {
	u, ok := m.urls[code]
	if !ok {
		return "", errors.New("url not found")
	}
	return u, nil
}

func TestShortenAndGetURL(t *testing.T) {
	m := newMockQuerier()
	s := NewService(m)
	ctx := context.Background()

	originalURL := "https://example.com/some/long/path"
	code, err := s.Shorten(ctx, originalURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(code) != LENGTH {
		t.Errorf("expected code length %d, got %d (code: %s)", LENGTH, len(code), code)
	}

	gotURL, err := s.GetURL(ctx, code)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotURL != originalURL {
		t.Errorf("expected %s, got %s", originalURL, gotURL)
	}

	nonExistent, err := s.GetURL(ctx, "invalid")
	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
	if nonExistent != "" {
		t.Errorf("expected empty string for non-existent code, got %s", nonExistent)
	}
}

func TestGenerateCode(t *testing.T) {
	for range 100 {
		code, err := generateCode()
		if err != nil {
			t.Fatalf("generateCode failed: %v", err)
		}

		if len(code) != LENGTH {
			t.Errorf("expected code length %d, got %d", LENGTH, len(code))
		}

		for _, char := range code {
			if !strings.ContainsRune(CHARSET, char) {
				t.Errorf("character %q not in CHARSET", char)
			}
		}
	}
}
