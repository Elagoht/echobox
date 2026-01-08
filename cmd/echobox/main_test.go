package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {
	// Test with default config
	server := createServer()

	if server == nil {
		t.Fatal("createServer() returned nil")
	}

	// Check that server is configured
	if server.Handler == nil {
		t.Error("createServer() handler is nil")
	}

	if server.ReadTimeout == 0 {
		t.Error("createServer() ReadTimeout is 0")
	}

	if server.WriteTimeout == 0 {
		t.Error("createServer() WriteTimeout is 0")
	}

	// Test with custom config
	oldPort := os.Getenv("PORT")
	oldRead := os.Getenv("READ_TIMEOUT")
	oldWrite := os.Getenv("WRITE_TIMEOUT")

	defer func() {
		os.Setenv("PORT", oldPort)
		os.Setenv("READ_TIMEOUT", oldRead)
		os.Setenv("WRITE_TIMEOUT", oldWrite)
	}()

	os.Setenv("PORT", "9999")
	os.Setenv("READ_TIMEOUT", "10")
	os.Setenv("WRITE_TIMEOUT", "20")

	server = createServer()

	if server.Addr != ":9999" {
		t.Errorf("createServer() Addr = %v, want :9999", server.Addr)
	}

	expectedReadTimeout := 10 * time.Second
	if server.ReadTimeout != expectedReadTimeout {
		t.Errorf("createServer() ReadTimeout = %v, want %v", server.ReadTimeout, expectedReadTimeout)
	}

	expectedWriteTimeout := 20 * time.Second
	if server.WriteTimeout != expectedWriteTimeout {
		t.Errorf("createServer() WriteTimeout = %v, want %v", server.WriteTimeout, expectedWriteTimeout)
	}
}

func TestCreateServerHandler(t *testing.T) {
	server := createServer()

	// Test that the handler works
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Handler returned status %v, want %v", w.Code, http.StatusOK)
	}
}

func TestRunServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping server test in short mode")
	}

	// Create a test server on a random port
	oldPort := os.Getenv("PORT")
	defer os.Setenv("PORT", oldPort)

	os.Setenv("PORT", "5872")
	server := createServer()

	// Start server in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- runServer(ctx, server)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that server is running
	resp, err := http.Get("http://localhost:5872/")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Shutdown server
	cancel()

	// Wait for server to stop (with timeout)
	select {
	case err := <-errChan:
		// Server stopped
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server stopped with error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Server did not stop within timeout")
	}
}

func TestRunServer_StartError(t *testing.T) {
	// Test error path when server fails to start
	// We'll create a server with an invalid address
	server := &http.Server{
		Addr:    ":invalid",
		Handler: createServer().Handler,
	}

	ctx := context.Background()
	err := runServer(ctx, server)

	// Should return an error (not panic)
	if err == nil {
		t.Error("Expected error when starting server with invalid address, got nil")
	}
}

func TestRunServer_Shutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping shutdown test in short mode")
	}

	oldPort := os.Getenv("PORT")
	defer os.Setenv("PORT", oldPort)

	os.Setenv("PORT", "5874")
	server := createServer()

	// Start server and immediately shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		errChan <- runServer(ctx, server)
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Trigger shutdown
	cancel()

	// Wait for server to stop
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected error from runServer: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Server did not stop within timeout")
	}
}

func TestMainFunc_StartError(t *testing.T) {
	// Test main() when server fails to start
	// We'll use an invalid port to cause startup failure
	if os.Getenv("TEST_MAIN_ERROR") == "1" {
		os.Setenv("PORT", "invalid")
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainFunc_StartError")
	cmd.Env = append(os.Environ(), "TEST_MAIN_ERROR=1")
	output, err := cmd.CombinedOutput()

	// Should exit with non-zero status
	if err == nil {
		t.Error("Expected error when main() fails to start server, got nil")
	}

	// Should contain error message in output
	if len(output) > 0 {
		t.Logf("main() output: %s", output)
	}
}

func TestInitLogging(t *testing.T) {
	// Test that init() runs without panicking
	// The init function is called automatically before tests run
	// If we reach here, init() succeeded
	if os.Getenv("TEST_INIT_LOGGING") == "1" {
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitLogging")
	cmd.Env = append(os.Environ(), "TEST_INIT_LOGGING=1")
	if err := cmd.Run(); err != nil {
		t.Fatalf("TestInitLogging failed: %v", err)
	}
}

func TestIntegration_AllEndpoints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Save original env
	oldPort := os.Getenv("PORT")

	defer func() {
		os.Setenv("PORT", oldPort)
	}()

	os.Setenv("PORT", "5873")

	// Start server in goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := createServer()
	errChan := make(chan error, 1)
	go func() {
		errChan <- runServer(ctx, server)
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	defer func() {
		cancel()
		// Wait for server shutdown (don't fail test if it takes too long)
		select {
		case <-errChan:
		case <-time.After(2 * time.Second):
		}
	}()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{"GET root", "GET", "/", "", http.StatusOK},
		{"POST root", "POST", "/", `{"test":"data"}`, http.StatusOK},
		{"GET headers", "GET", "/headers", "", http.StatusOK},
		{"POST body", "POST", "/body", "test body", http.StatusOK},
		{"GET queries", "GET", "/queries?foo=bar", "", http.StatusOK},
		{"GET 200", "GET", "/200", "", http.StatusOK},
		{"GET 404", "GET", "/404", "", http.StatusNotFound},
		{"GET 500", "GET", "/500", "", http.StatusInternalServerError},
		{"PUT root", "PUT", "/", "data", http.StatusOK},
		{"DELETE root", "DELETE", "/", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			}

			req, err := http.NewRequest(tt.method, "http://localhost:5873"+tt.path, bodyReader)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %v, got %v", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
