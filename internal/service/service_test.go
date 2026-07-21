package service

import (
	"strings"
	"sync"
	"testing"
)

func TestShortenAndGetURL(t *testing.T) {
	s := NewService()

	originalURL := "https://example.com/some/long/path"
	code, err := s.Shorten(originalURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(code) != LENGTH {
		t.Errorf("expected code length %d, got %d (code: %s)", LENGTH, len(code), code)
	}

	gotURL := s.GetURL(code)
	if gotURL != originalURL {
		t.Errorf("expected %s, got %s", originalURL, gotURL)
	}

	nonExistent := s.GetURL("invalid")
	if nonExistent != "" {
		t.Errorf("expected empty string for non-existent code, got %s", nonExistent)
	}
}

func TestConcurrentShortenAndGet(t *testing.T) {
	s := NewService()
	var wg sync.WaitGroup
	count := 50

	// Concurrently shorten URLs
	codes := make([]string, count)
	for i := range count {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			code, err := s.Shorten("https://example.com")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			codes[idx] = code
		}(i)
	}
	wg.Wait()

	// Concurrently read URLs
	for _, code := range codes {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			url := s.GetURL(c)
			if url != "https://example.com" {
				t.Errorf("expected https://example.com, got %s", url)
			}
		}(code)
	}
	wg.Wait()
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
