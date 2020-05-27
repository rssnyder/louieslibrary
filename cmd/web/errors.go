package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

// ServerError logs stack and send server error
func (app *App) ServerError(w http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// ClientError sends error in response to client
func (app *App) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// NotFound is generic not found error
func (app *App) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
