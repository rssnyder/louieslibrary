package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Mr-Schneider/louieslibrary/pkg/models"
)

// HTMLData models the page data
type HTMLData struct {
	Request     *models.Request
	Requests    []*models.Request
	User        *models.User
	DisplayUser *models.User
	Book        *models.Book
	Books       []*models.Book
	Reviews     []*models.Review
	Invites     []*models.Invite
	Messages    []*models.Message
	Threads     []*models.Message
	Path        string
	Form        interface{}
	Flash       string
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

	// Check logged in status
	user := &models.User{}
	_, user = app.LoggedIn(r)
	data.User = user

	// Get unread messages
	unread, err := app.DB.GetUnopened(user.Username)
	if len(unread) != 0 {
		session.AddFlash("You have new messages!", "default")
	}

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

	// Write template to buffer, then send buffer
	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}

// JSONResponse sends a response in json format
func JSONResponse(w http.ResponseWriter, code int, output interface{}) {

	// Convert our interface to JSON
	response, _ := json.Marshal(output)

	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")

	// Our response code
	w.WriteHeader(code)

	w.Write(response)
}
