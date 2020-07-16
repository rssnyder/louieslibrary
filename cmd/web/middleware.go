package main

import (
	"log"
	"net/http"
	"strings"
)

// LogRequest Logs every request
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// SecureHeaders Set secure headers
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}

		next.ServeHTTP(w, r)
	})
}

// RequireLogin 401 for unathed users
func (app *App) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header["Authorization"]
		if len(token) > 0 {
			splits := strings.Split(token[0], " ")
			if len(splits) > 0 {
				_, _, _, err := app.VerifyJWT(splits[1])
				if err != nil {
					log.Println("Issues verifying token")
				} else {
					next.ServeHTTP(w, r)
					return
				}
			} else {
				log.Println("Invalid token passed in header")
			}
		} else {
			log.Println("No authorization header")
		}

		JSONResponse(w, 401, "")
	})
}

// RequireWriter Redirect users without writer role
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
