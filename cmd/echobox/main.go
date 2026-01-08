package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Elagoht/echobox/internal/config"
	"github.com/Elagoht/echobox/internal/router"
)

func main() {
	cfg := config.Load()

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.New(),
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	log.Printf("Echobox listening on http://localhost:%s", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("[%s] ", "echobox"))
}
