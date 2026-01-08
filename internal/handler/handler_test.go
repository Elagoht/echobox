package handler

import (
	"errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

type errorWriter struct {
	http.ResponseWriter
	writeErr bool
}

func (e *errorWriter) Write(b []byte) (int, error) {
	if e.writeErr {
		return 0, errors.New("write error")
	}
	return e.ResponseWriter.Write(b)
}

type brokenWriter struct {
	http.ResponseWriter
}

func (b *brokenWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func (b *brokenWriter) WriteHeader(int) {}


func TestEcho(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		headers    map[string]string
		query      string
		wantStatus int
	}{
		{
			name:       "simple GET request",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST with body",
			method:     http.MethodPost,
			body:       `{"test":"data"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "PUT with body",
			method:     http.MethodPut,
			body:       "test body",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/?"+tt.query, strings.NewReader(tt.body))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			Echo(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Echo() status = %v, want %v", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Echo() Content-Type = %v, want application/json", contentType)
			}

			var resp EchoResponse
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if resp.Method != tt.method {
				t.Errorf("Echo() method = %v, want %v", resp.Method, tt.method)
			}

			if resp.Body != tt.body {
				t.Errorf("Echo() body = %v, want %v", resp.Body, tt.body)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-Header", "test-value")
	req.Header.Set("User-Agent", "test-agent")

	w := httptest.NewRecorder()
	Headers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Headers() status = %v, want %v", w.Code, http.StatusOK)
	}

	var headers map[string][]string
	if err := json.NewDecoder(w.Body).Decode(&headers); err != nil {
		t.Fatalf("Failed to decode headers: %v", err)
	}

	if headers["X-Custom-Header"][0] != "test-value" {
		t.Errorf("Headers() X-Custom-Header = %v, want test-value", headers["X-Custom-Header"])
	}
}

func TestBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "json body",
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello"}`,
		},
		{
			name:       "plain text",
			body:       "plain text body",
			wantStatus: http.StatusOK,
			wantBody:   "plain text body",
		},
		{
			name:       "empty body",
			body:       "",
			wantStatus: http.StatusOK,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			Body(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Body() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if w.Body.String() != tt.wantBody {
				t.Errorf("Body() body = %v, want %v", w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestQueries(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?foo=bar&baz=qux&foo=second", nil)
	w := httptest.NewRecorder()

	Queries(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Queries() status = %v, want %v", w.Code, http.StatusOK)
	}

	var queries map[string][]string
	if err := json.NewDecoder(w.Body).Decode(&queries); err != nil {
		t.Fatalf("Failed to decode queries: %v", err)
	}

	if queries["foo"][0] != "bar" {
		t.Errorf("Queries() foo[0] = %v, want bar", queries["foo"][0])
	}
	if queries["baz"][0] != "qux" {
		t.Errorf("Queries() baz[0] = %v, want qux", queries["baz"][0])
	}
	if len(queries["foo"]) != 2 {
		t.Errorf("Queries() foo length = %v, want 2", len(queries["foo"]))
	}
}

func TestMatchStatusCode(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/200", true},
		{"/404", true},
		{"/500", true},
		{"/699", true},
		{"/199", false},
		{"/700", false},
		{"/2000", false},
		{"/20", false},
		{"/", false},
		{"/abc", false},
		{"/headers", false},
		{"/body", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := MatchStatusCode(tt.path)
			if got != tt.want {
				t.Errorf("MatchStatusCode(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestServeStatusCode(t *testing.T) {
	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/200", 200, "OK"},
		{"/404", 404, "Not Found"},
		{"/500", 500, "Internal Server Error"},
		{"/418", 418, "I'm a teapot"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			ServeStatusCode(w, tt.path)

			if w.Code != tt.wantStatus {
				t.Errorf("ServeStatusCode() status = %v, want %v", w.Code, tt.wantStatus)
			}

			body := strings.TrimSpace(w.Body.String())
			if body != tt.wantBody {
				t.Errorf("ServeStatusCode() body = %v, want %v", body, tt.wantBody)
			}
		})
	}
}

func TestMethodAllow(t *testing.T) {
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrapped := MethodAllow(testHandler)

	tests := []struct {
		method string
	}{
		{http.MethodGet},
		{http.MethodPost},
		{http.MethodPut},
		{http.MethodDelete},
		{http.MethodPatch},
		{http.MethodHead},
		{http.MethodOptions},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			handlerCalled = false
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			wrapped(w, req)

			if !handlerCalled {
				t.Error("MethodAllow() handler was not called")
			}

			allow := w.Header().Get("Allow")
			if allow != "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS" {
				t.Errorf("MethodAllow() Allow header = %v, want GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS", allow)
			}
		})
	}
}

func TestEcho_ErrorReadingBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", errorReader{})
	w := httptest.NewRecorder()

	Echo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Echo() error status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestEcho_ErrorEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Use brokenWriter to cause encoding errors
	bw := &brokenWriter{ResponseWriter: w}
	Echo(bw, req)

	// Should not panic, just log the error
	_ = w.Code
}

func TestHeaders_ErrorEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	bw := &brokenWriter{ResponseWriter: w}
	Headers(bw, req)

	// Should not panic
	_ = w.Code
}

func TestBody_ErrorReading(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", errorReader{})
	w := httptest.NewRecorder()

	Body(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Body() error status = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestBody_ErrorWriting(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("test"))
	w := httptest.NewRecorder()

	bw := &brokenWriter{ResponseWriter: w}
	Body(bw, req)

	// Should not panic
	_ = w.Code
}

func TestQueries_ErrorEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?foo=bar", nil)
	w := httptest.NewRecorder()

	bw := &brokenWriter{ResponseWriter: w}
	Queries(bw, req)

	// Should not panic
	_ = w.Code
}

func TestServeStatusCode_UnknownCode(t *testing.T) {
	w := httptest.NewRecorder()

	// Test with a non-standard status code
	ServeStatusCode(w, "/699")

	if w.Code != 699 {
		t.Errorf("ServeStatusCode() status = %v, want 699", w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	if body != "Unknown Status Code" {
		t.Errorf("ServeStatusCode() body = %v, want 'Unknown Status Code'", body)
	}
}

func TestServeStatusCode_ErrorWriting(t *testing.T) {
	w := httptest.NewRecorder()
	bw := &brokenWriter{ResponseWriter: w}

	ServeStatusCode(bw, "/404")

	// Should not panic
	_ = w.Code
}

