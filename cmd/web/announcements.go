package main

import (
	"log"
	"net/http"

	"github.com/rssnyder/louieslibrary/pkg/forms"
)

// NewAnnouncement display the new announcement form
func (app *App) NewAnnouncement(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newannouncement.page.html", &HTMLData{
		Form: &forms.NewAnnouncement{},
	})
}

// CreateAnnouncement set a new announcement
func (app *App) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {

	// Get current user
	_, user := app.LoggedIn(r)

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new announcement on the information from the form
	announcement := &forms.NewAnnouncement{
		Author:  user.Username,
		Content: r.PostForm.Get("content"),
	}

	// Validate the new announcement form
	// if !announcement.Valid() {
	// 	app.RenderHTML(w, r, "newannouncement.page.html", &HTMLData{Form: announcement})
	// 	return
	// }

	log.Printf("New announcement by %s: %s", announcement.Author, announcement.Content)

	// Insert the new announcement
	_, err = app.DB.InsertAnnouncement(announcement)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Direct to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
