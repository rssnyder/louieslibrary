package main

import (
	"net/http"
)

// Home page of site
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	// 404 if not truly root
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	// Get the latest requests
	requests, err := app.DB.LatestRequests()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get the latest books
	books, err := app.DB.LatestBooks()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Requests: requests,
		Books:    books,
	})
}