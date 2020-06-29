package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/forms"
	"github.com/gorilla/mux"
)

// ShowRequest displays a single request
func (app *App) ShowRequest(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get requested snippet id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	// Get request
	request, err := app.DB.GetRequest(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if request == nil {
		app.NotFound(w)
		return
	}

	request.BookID = strings.TrimSpace(request.BookID)

	// Get the previous flashes, if any.
	if flashes := session.Flashes("default"); len(flashes) > 0 {
		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, "showrequest.page.html", &HTMLData{
			Request: request,
			Flash:   fmt.Sprintf("%v", flashes[0]),
		})
	} else {
		app.RenderHTML(w, r, "showrequest.page.html", &HTMLData{
			Request: request,
			Flash:   "",
		})
	}
}

// NewRequest displays the new request form
func (app *App) NewRequest(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newrequest.page.html", &HTMLData{
		Form: &forms.NewRequest{},
	})
}

// CreateRequest creates a new request
func (app *App) CreateRequest(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Get requester
	valid, user := app.LoggedIn(r)
	if !valid {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new request based on html form
	form := &forms.NewRequest{
		Requester: user.Username,
		Title:     r.PostForm.Get("title"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "newrequest.page.html", &HTMLData{Form: form})
		return
	}

	// Insert the new request
	id, err := app.DB.InsertRequest(form.Requester, form.Title, r.RemoteAddr)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Your request was saved successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/request/%d", id), http.StatusSeeOther)
}

// FillRequest fills q request
func (app *App) FillRequest(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get requested snippet id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	// Parse the post data
	err = r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Insert the new request
	reply := app.DB.FillRequest(id, r.PostForm.Get("bookid"))
	if reply == "" {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Request filled!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/request/%d", id), http.StatusSeeOther)
}

// ListAllRequests does what it says
func (app *App) ListAllRequests(w http.ResponseWriter, r *http.Request) {
	// Get the requests
	requests, err := app.DB.LatestRequests(1000)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "showrequests.page.html", &HTMLData{
		Requests:    requests,
	})
}