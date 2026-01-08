package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func Echo(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp := EchoResponse{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    string(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func Headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(r.Header); err != nil {
		log.Printf("Error encoding headers: %v", err)
	}
}

func Body(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("Error writing body: %v", err)
	}
}

func Queries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(r.URL.Query()); err != nil {
		log.Printf("Error encoding queries: %v", err)
	}
}

func MatchStatusCode(path string) bool {
	matched, _ := regexp.MatchString("^/\\d{3}$", path)
	if !matched {
		return false
	}

	code, err := strconv.Atoi(path[1:])
	return err == nil && code >= 200 && code <= 699
}

func ServeStatusCode(w http.ResponseWriter, path string) {
	code, _ := strconv.Atoi(path[1:])

	w.WriteHeader(code)
	statusText := http.StatusText(code)
	if statusText == "" {
		statusText = "Unknown Status Code"
	}
	if _, err := w.Write([]byte(statusText)); err != nil {
		log.Printf("Error writing status: %v", err)
	}
}

func MethodAllow(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
		h(w, r)
	}
}
