package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	mux := New()

	if mux == nil {
		t.Fatal("New() returned nil")
	}
}

func TestRouter_Routes(t *testing.T) {
	mux := New()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantAllow  string
	}{
		{
			name:       "GET /headers",
			method:     http.MethodGet,
			path:       "/headers",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "POST /body",
			method:     http.MethodPost,
			path:       "/body",
			body:       "test body",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET /queries",
			method:     http.MethodGet,
			path:       "/queries",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET / (echo)",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET /404 (status code)",
			method:     http.MethodGet,
			path:       "/404",
			wantStatus: http.StatusNotFound,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET /500 (status code)",
			method:     http.MethodGet,
			path:       "/500",
			wantStatus: http.StatusInternalServerError,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET /200 (status code)",
			method:     http.MethodGet,
			path:       "/200",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
		{
			name:       "GET /randompath (echo)",
			method:     http.MethodGet,
			path:       "/randompath",
			wantStatus: http.StatusOK,
			wantAllow:  "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Router status = %v, want %v", w.Code, tt.wantStatus)
			}

			allow := w.Header().Get("Allow")
			if allow != tt.wantAllow {
				t.Errorf("Router Allow header = %v, want %v", allow, tt.wantAllow)
			}
		})
	}
}

func TestRouter_AllowsAllMethods(t *testing.T) {
	mux := New()

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/headers", nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			allow := w.Header().Get("Allow")
			expectedAllow := "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS"
			if allow != expectedAllow {
				t.Errorf("%s: Allow header = %v, want %v", method, allow, expectedAllow)
			}
		})
	}
}

func TestRouter_StatusCodeRange(t *testing.T) {
	mux := New()

	statusCodes := []string{"/200", "/201", "/204", "/301", "/302", "/400", "/401", "/403", "/404", "/500", "/502", "/503"}

	for _, code := range statusCodes {
		t.Run(code, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, code, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			// Just verify we get a response (not 404 from the router itself)
			// The status code handler should return the requested code
			if w.Code == http.StatusNotFound && code != "/404" {
				t.Errorf("%s: got 404, expected status code to be returned", code)
			}
		})
	}
}
