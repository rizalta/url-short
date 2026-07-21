package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	urls map[string]string
}

func newMockService() *mockService {
	return &mockService{
		urls: make(map[string]string),
	}
}

func (m *mockService) Shorten(url string) (string, error) {
	code := "abc123"
	m.urls[code] = url
	return code, nil
}

func (m *mockService) GetURL(code string) string {
	return m.urls[code]
}

func TestShortenHandler(t *testing.T) {
	mockSvc := newMockService()
	h := NewHandler(mockSvc)
	router := h.Routes()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		checkBody      bool
	}{
		{
			name:           "Valid URL",
			body:           `{"url": "https://example.com"}`,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "URL missing scheme gets auto-prefixed",
			body:           `{"url": "example.com"}`,
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "Invalid JSON",
			body:           `{invalid-json}`,
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "Empty URL",
			body:           `{"url": ""}`,
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkBody {
				var res ShortenRes
				if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if res.ShortURL == "" {
					t.Errorf("expected non-empty short_url")
				}
			}
		})
	}
}

func TestGetURLHandler(t *testing.T) {
	mockSvc := newMockService()
	mockSvc.urls["abc123"] = "https://example.com"

	h := NewHandler(mockSvc)
	router := h.Routes()

	t.Run("Existing Code Redirects", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, w.Code)
		}

		location := w.Header().Get("Location")
		if location != "https://example.com" {
			t.Errorf("expected Location header 'https://example.com', got %q", location)
		}
	})

	t.Run("Non-existent Code Returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
