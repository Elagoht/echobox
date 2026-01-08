package router

import (
	"github.com/Elagoht/echobox/internal/handler"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// Apply method allow middleware to all handlers
	mux.HandleFunc("/headers", handler.MethodAllow(handler.Headers))
	mux.HandleFunc("/body", handler.MethodAllow(handler.Body))
	mux.HandleFunc("/queries", handler.MethodAllow(handler.Queries))

	// Catch-all handler for status codes and echo
	mux.HandleFunc("/", handler.MethodAllow(func(w http.ResponseWriter, r *http.Request) {
		// Check if path is a 3-digit status code
		if handler.MatchStatusCode(r.URL.Path) {
			handler.ServeStatusCode(w, r.URL.Path)
			return
		}

		// Default to echo handler
		handler.Echo(w, r)
	}))

	return mux
}
