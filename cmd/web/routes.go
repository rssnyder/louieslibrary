package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes defines site routes
func (app *App) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", app.Home).Methods("GET")

	// Requests
	r.Handle("/request/new", app.RequireLogin(http.HandlerFunc(app.NewRequest))).Methods("GET")
	r.Handle("/request/new", app.RequireLogin(http.HandlerFunc(app.CreateRequest))).Methods("POST")
	r.Handle("/request/{id}", app.RequireLogin(http.HandlerFunc(app.ShowRequest))).Methods("GET")

	// Fill requests
	r.Handle("/request/{id}/fill", app.RequireWriter(http.HandlerFunc(app.FillRequest))).Methods("POST")

	// Books
	r.Handle("/book/all", app.RequireLogin(http.HandlerFunc(app.ListAllBooks))).Methods("GET")
	r.Handle("/book/{id}", app.RequireLogin(http.HandlerFunc(app.ShowBook))).Methods("GET")
	r.Handle("/book/review", app.RequireLogin(http.HandlerFunc(app.CreateReview))).Methods("POST")

	// Adding books
	r.Handle("/write/book", app.RequireWriter(http.HandlerFunc(app.NewBook))).Methods("GET")
	r.Handle("/write/book", app.RequireWriter(http.HandlerFunc(app.CreateBook))).Methods("POST")

	// Youtube
	r.Handle("/youtube/playlist", app.RequireLogin(http.HandlerFunc(app.NewPlaylist))).Methods("GET")
	r.Handle("/youtube/playlist", app.RequireLogin(http.HandlerFunc(app.DownloadPlaylist))).Methods("POST")

	// Unlock user methods
	r.HandleFunc("/user/signup", app.SignupUser).Methods("GET")
	r.HandleFunc("/user/signup", app.CreateUser).Methods("POST")
	r.HandleFunc("/user/login", app.LoginUser).Methods("GET")
	r.HandleFunc("/user/login", app.VerifyUser).Methods("POST")
	r.HandleFunc("/user/logout", app.LogoutUser).Methods("GET")
	r.Handle("/user/{username}", app.RequireLogin(http.HandlerFunc(app.ShowUser))).Methods("GET")

	// Book files
	bookServer := http.FileServer(http.Dir(app.BookDir))
	r.PathPrefix("/books/").Handler(http.StripPrefix("/books/", bookServer))

	// Youtube files
	ytServer := http.FileServer(http.Dir(app.YoutubeDir))
	r.PathPrefix("/youtube/").Handler(http.StripPrefix("/youtube/", ytServer))

	// Fileserver for css and js files
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.Use(SecureHeaders)
	r.Use(LogRequest)

	return r
}
