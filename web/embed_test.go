package web

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticMiddleware(t *testing.T) {
	staticMiddleware, err := StaticMiddleware()
	if err != nil {
		t.Fatalf("failed to initialize static middleware: %v", err)
	}

	dist, err := fs.Sub(distFS, "dist")
	if err != nil {
		t.Fatalf("failed to get sub fs: %v", err)
	}

	entries, err := fs.ReadDir(dist, "assets")
	if err != nil || len(entries) == 0 {
		t.Fatalf("failed to get assets: %v", err)
	}

	assetPath := "/assets/" + entries[0].Name()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	handler := staticMiddleware(nextHandler)

	tests := []struct {
		name           string
		requestPath    string
		expectedStatus int
	}{
		{
			name:           "serves index.html on root path",
			requestPath:    "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "serves static asset favicon",
			requestPath:    "/favicon.svg",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "serves assets",
			requestPath:    assetPath,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "falls through next handler when not static route",
			requestPath:    "/abcd",
			expectedStatus: http.StatusTeapot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
