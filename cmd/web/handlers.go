package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"request.thecornelius.duckdns.org/pkg/forms"
)

// Home page of site
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	// 404 if not truly root
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	// Get the latest requests
	requests, err := app.Database.LatestRequests()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Requests: requests,
	})
}

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
	request, err := app.Database.GetRequest(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if request == nil {
		app.NotFound(w)
		return
	}

	// Get the previous flashes, if any.
	if flashes := session.Flashes("default"); len(flashes) > 0 {
		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, "show.page.html", &HTMLData{
			Request: request,
			Flash:   fmt.Sprintf("%v", flashes[0]),
		})
	} else {
		app.RenderHTML(w, r, "show.page.html", &HTMLData{
			Request: request,
			Flash:   "",
		})
	}
}

// NewRequest displays the new request form
func (app *App) NewRequest(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "new.page.html", &HTMLData{
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

	// Model the new request based on html form
	form := &forms.NewRequest{
		Requester: r.PostForm.Get("requester"),
		Title:     r.PostForm.Get("title"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "new.page.html", &HTMLData{Form: form})
		return
	}

	// Insert the new request
	id, err := app.Database.InsertRequest(form.Requester, form.Title, r.RemoteAddr)
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
