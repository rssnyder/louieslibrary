package main

import (
	"log"
	"net/http"
)

// LogRequest logs every request on the server
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders sets secure headers on every request
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}

		next.ServeHTTP(w, r)
	})
}
