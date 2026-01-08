package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Elagoht/echobox/internal/config"
	"github.com/Elagoht/echobox/internal/router"
)

func createServer() *http.Server {
	cfg := config.Load()

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.New(),
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}
}

func runServer(ctx context.Context, server *http.Server) error {
	log.Printf("Echobox listening on http://localhost:%s", server.Addr[1:])

	// Channel to capture server errors
	errChan := make(chan error, 1)

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server failed to start: %w", err)
		} else {
			errChan <- nil
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		// Context cancelled, shutdown server gracefully
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
		// Return any error from ListenAndServe
		return <-errChan
	case err := <-errChan:
		// Server stopped (either error or closed)
		return err
	}
}

func main() {
	server := createServer()

	if err := runServer(context.Background(), server); err != nil {
		log.Fatalf("%v", err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("[%s] ", "echobox"))
}
