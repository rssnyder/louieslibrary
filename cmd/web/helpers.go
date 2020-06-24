package main

import (
	"net/http"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
)

// LoggedIn gets the logged in status, also returns user object
func (app *App) LoggedIn(r *http.Request) (bool, *models.User) {
	session, _ := app.Sessions.Get(r, "session-name")

	// Get user from session
	loggedIn := session.Values["user"]
	var user = &models.User{}
	user, ok := loggedIn.(*models.User)

	// Test if user exists
	if !ok {
		return false, &models.User{}
	}

	// Check if user not empty
	if user == (&models.User{}) {
		return false, &models.User{}
	}

	// Check if user not empty still
	if user.ID == 0 {
		return false, &models.User{}
	}

	// User is logged in
	return true, user
}