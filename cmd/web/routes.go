package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes define the site routes
func (app *App) Routes() *mux.Router {

	// Create a new router
	r := mux.NewRouter()

	// Homepage
	r.Handle("/", app.RequireLogin(http.HandlerFunc(app.Home))).Methods("GET")
	r.Handle("/about", app.RequireLogin(http.HandlerFunc(app.About))).Methods("GET")

	// Requests
	r.Handle("/request/all", app.RequireLogin(http.HandlerFunc(app.ListAllRequests))).Methods("GET")
	r.Handle("/request/new", app.RequireLogin(http.HandlerFunc(app.NewRequest))).Methods("GET")
	r.Handle("/request/new", app.RequireLogin(http.HandlerFunc(app.CreateRequest))).Methods("POST")
	r.Handle("/request/{id}", app.RequireLogin(http.HandlerFunc(app.ShowRequest))).Methods("GET")
	r.Handle("/request/{id}/fill", app.RequireWriter(http.HandlerFunc(app.FillRequest))).Methods("POST")

	// Books
	r.Handle("/book/all", app.RequireLogin(http.HandlerFunc(app.ListAllBooks))).Methods("GET")
	r.Handle("/book/review", app.RequireLogin(http.HandlerFunc(app.CreateReview))).Methods("POST")
	r.Handle("/book/edit", app.RequireLogin(http.HandlerFunc(app.UpdateBook))).Methods("POST")
	r.Handle("/book/edit/{volumeid}", app.RequireLogin(http.HandlerFunc(app.EditBook))).Methods("GET")
	r.Handle("/book/collect/{volumeid}", app.RequireLogin(http.HandlerFunc(app.AddToCollection))).Methods("POST")
	r.Handle("/book/{volumeid}", app.RequireLogin(http.HandlerFunc(app.ShowBook))).Methods("GET")
	r.Handle("/book/{volumeid}", app.RequireLogin(http.HandlerFunc(app.DownloadBook))).Methods("POST")
	r.Handle("/write/book", app.RequireWriter(http.HandlerFunc(app.NewBook))).Methods("GET")
	r.Handle("/write/book", app.RequireWriter(http.HandlerFunc(app.CreateBook))).Methods("POST")

	// Messages
	r.Handle("/messages/{reciver}", app.RequireLogin(http.HandlerFunc(app.Messages))).Methods("GET")
	r.Handle("/messages/{reciver}", app.RequireLogin(http.HandlerFunc(app.CreateMessage))).Methods("POST")

	// Announcements
	r.Handle("/announcement/new", app.RequireWriter(http.HandlerFunc(app.NewAnnouncement))).Methods("GET")
	r.Handle("/announcement/new", app.RequireWriter(http.HandlerFunc(app.CreateAnnouncement))).Methods("POST")

	// Youtube
	r.Handle("/youtube/playlist", app.RequireLogin(http.HandlerFunc(app.NewPlaylist))).Methods("GET")
	r.Handle("/youtube/playlist", app.RequireLogin(http.HandlerFunc(app.DownloadPlaylist))).Methods("POST")

	// Unlocked user methods
	r.HandleFunc("/user/signup", app.SignupUser).Methods("GET")
	r.HandleFunc("/user/signup", app.CreateUser).Methods("POST")
	r.HandleFunc("/user/login", app.LoginUser).Methods("GET")
	r.HandleFunc("/user/login", app.VerifyUser).Methods("POST")
	r.HandleFunc("/user/logout", app.LogoutUser).Methods("GET")
	r.Handle("/user/invite/create", app.RequireLogin(http.HandlerFunc(app.CreateInviteCode))).Methods("POST")
	r.Handle("/user/{username}", app.RequireLogin(http.HandlerFunc(app.ShowUser))).Methods("GET")

	// Hosting static files

	// Youtube files
	ytServer := http.FileServer(http.Dir(app.YoutubeDir))
	r.PathPrefix("/youtube/").Handler(http.StripPrefix("/youtube/", ytServer))

	// CSS and HTML files
	fileServer := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// Global middleware
	r.Use(SecureHeaders)
	r.Use(LogRequest)

	return r
}
