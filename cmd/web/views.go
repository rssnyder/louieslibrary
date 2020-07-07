package main

import (
	"fmt"
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
)

// HTMLData models the page data
type HTMLData struct {
	Request  		*models.Request
	Requests 		[]*models.Request
	User     		*models.User
	Users				[]*models.User
	DisplayUser	*models.User
	Book     		*models.Book
	Books    		[]*models.Book
	Reviews  		[]*models.Review
	Invites  		[]*models.Invite
	Messages		[]*models.Message
	Friends			[]*models.Message
	Path     		string
	Form     		interface{}
	Flash    		string
}

// humanDate
// Format dates in a better view
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// RenderHTML
// Display the current page based on htmldata
func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, page string, data *HTMLData) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// If no data provided, create emtpy data
	if data == nil {
		data = &HTMLData{}
	}

	// Add the current path to the data
	data.Path = r.URL.Path

	// Get the previous flash
	if flashes := session.Flashes("default"); len(flashes) > 0 {

		// Save session
		err := session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		data.Flash = fmt.Sprintf("%v", flashes[0])
	}

	// Check logged in status
	user := &models.User{}
	_, user = app.LoggedIn(r)
	data.User = user

	// Render the base template with target page
	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	// Map for custome template functions
	fm := template.FuncMap{
		"humanDate": humanDate,
	}

	// Pull the html files together
	ts, err := template.New("").Funcs(fm).ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Write template to buffer, then send buffer
	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}
