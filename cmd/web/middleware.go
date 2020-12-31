package main

import (
	"log"
	"net/http"
)

// LogRequest logs every request
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders set secure headers
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}

		next.ServeHTTP(w, r)
	})
}

// RequireLogin redirect unauthenticated users
func (app *App) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// WebUI
		loggedIn, _ := app.LoggedIn(r)
		if !loggedIn {

			// Try for a jwt
			if app.ValidateRequest(w, r) {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/user/login", 302)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireWriter redirect users without writer role
func (app *App) RequireWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedIn, user := app.LoggedIn(r)
		if !loggedIn {
			http.Redirect(w, r, "/user/login", 302)
			return
		}

		if user.Role != "writer" {
			http.Redirect(w, r, "/", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}
