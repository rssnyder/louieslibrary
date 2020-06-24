package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes defines site routes
func (app *App) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", app.Home).Methods("GET")

	request := r.PathPrefix("/request").Subrouter()
	request.HandleFunc("/new", app.NewRequest).Methods("GET")
	request.HandleFunc("/new", app.CreateRequest).Methods("POST")
	request.HandleFunc("/{id}", app.ShowRequest).Methods("GET")
	request.Use(app.RequireLogin)

	r.HandleFunc("/user/signup", app.SignupUser).Methods("GET")
	r.HandleFunc("/user/signup", app.CreateUser).Methods("POST")

	r.HandleFunc("/user/login", app.LoginUser).Methods("GET")
	r.HandleFunc("/user/login", app.VerifyUser).Methods("POST")
	r.HandleFunc("/user/logout", app.LogoutUser).Methods("GET")

	// Fileserver for css and js files
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.Use(SecureHeaders)
	r.Use(LogRequest)

	return r
}
