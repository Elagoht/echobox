package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type EchoResponse struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    string(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Header)
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	w.Write(bodyBytes)
}

func queriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.URL.Query())
}

func methodHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all methods
		w.Header().Set("Allow", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		handler(w, r)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check if path is a 3-digit number between 200-699
	if matched, _ := regexp.MatchString("^/\\d{3}$", path); matched {
		code, err := strconv.Atoi(path[1:])
		if err == nil && code >= 200 && code <= 699 {
			w.WriteHeader(code)
			w.Write([]byte(http.StatusText(code)))
			return
		}
	}

	// Default to echo handler
	echoHandler(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5867"
	}

	http.HandleFunc("/headers", methodHandler(headersHandler))
	http.HandleFunc("/body", methodHandler(bodyHandler))
	http.HandleFunc("/queries", methodHandler(queriesHandler))
	http.HandleFunc("/", methodHandler(mainHandler))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
