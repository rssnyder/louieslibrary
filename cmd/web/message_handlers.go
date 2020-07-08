package main

import (
	"fmt"
	"net/http"
	//"strconv"
	//"strings"
	"github.com/Mr-Schneider/louieslibrary/pkg/forms"
	"github.com/gorilla/mux"
)

// NewMessage
// Display the new request form
func (app *App) Messages(w http.ResponseWriter, r *http.Request) {

	// Get requested conversation user
	vars := mux.Vars(r)
	reciver := vars["reciver"]

	// Get user
	_, user := app.LoggedIn(r)

	messages, err := app.DB.GetConversation(user.Username, reciver)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get user from db
	display_user, err := app.DB.GetUser(reciver)
	if err != nil {
		app.NotFound(w)
		return
	}

	// Get existing threads
	threads, err := app.DB.GetThreads(user.Username)

	// Get unread messages
	unread, err := app.DB.GetUnopened(user.Username)

	// Flag threads with unread messages
	// for _, closed := range unread {
	// 	for _, thread := range threads {
	// 		if closed.Sender == thread.Sender {
	// 			thread.Read = false
	// 			break
	// 		} else {
	// 			thread.Read = true
	// 		}
	// 	}
	// }

	for _, thread := range threads {
		for _, closed := range unread {
			if closed.Sender == thread.Sender {
				thread.Read = true
				break
			}
		}
	}


	// Model the new message form based on html form
	form := &forms.NewMessage{
		Reciver:  reciver,
	}

	app.RenderHTML(w, r, "messages.page.html", &HTMLData{
		Messages: messages,
		DisplayUser: display_user,
		Threads: threads,
		Form: form,
	})
}

// CreateMessage
// Create a new message in the db
func (app *App) CreateMessage(w http.ResponseWriter, r *http.Request) {
	
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Get sender
	_, user := app.LoggedIn(r)

	// Model the new message based on html form
	form := &forms.NewMessage{
		Sender: 	user.Username,
		Reciver:  r.PostForm.Get("reciver"),
		Content:	r.PostForm.Get("content"),
	}

	// Validate form
	// if !form.Valid() {
	// 	app.RenderHTML(w, r, "newrequest.page.html", &HTMLData{Form: form})
	// 	return
	// }

	// Insert the new request
	err = app.DB.InsertMessage(form.Sender, form.Reciver, form.Content)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session.AddFlash("Your message was saved successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	/// Send user to the newly added request
	http.Redirect(w, r, fmt.Sprintf("/messages/%s", form.Reciver), http.StatusSeeOther)
}
