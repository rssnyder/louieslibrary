package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes defines site routes
func (app *App) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", app.Home).Methods("GET")
	r.HandleFunc("/request/new", app.NewRequest).Methods("GET")
	r.HandleFunc("/request/new", app.CreateRequest).Methods("POST")
	r.HandleFunc("/request/{id}", app.ShowRequest).Methods("GET")

	// Fileserver for css and js files
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return r
}
