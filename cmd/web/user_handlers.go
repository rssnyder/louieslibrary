package main

import (
	"fmt"
	"net/http"

	"github.com/Mr-Schneider/louieslibrary/pkg/forms"
	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/gorilla/mux"
)

// SignupUser display the signup form
func (app *App) SignupUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.page.html", &HTMLData{
		Form: &forms.NewUser{},
	})
}

// CreateUser use signup form to create a new user
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new user based on html form
	form := &forms.NewUser{
		Username:   r.PostForm.Get("username"),
		Email:      r.PostForm.Get("email"),
		InviteCode: r.PostForm.Get("invitecode"),
		Password:   r.PostForm.Get("password"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	}

	// Validate request
	used, err := app.DB.ValidateInvite(form.InviteCode)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if used {
		// Save failure message
		session.AddFlash("Invalid invite code.", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Insert the new user
	err = app.DB.InsertUser(form.Username, form.Email, form.Password)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Fill invite
	err = app.DB.FillInvite(form.Username, form.InviteCode)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session.AddFlash("Your account was created successfully! Please login.", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// LoginUser displays a login page
func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get the previous flash
	if flashes := session.Flashes("default"); len(flashes) > 0 {

		// Save session
		err := session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Display login with flash
		app.RenderHTML(w, r, "login.page.html", &HTMLData{
			Form:  &forms.NewUser{},
			Flash: fmt.Sprintf("%v", flashes[0]),
		})
	} else {

		// Display login without flash
		app.RenderHTML(w, r, "login.page.html", &HTMLData{
			Form:  &forms.NewUser{},
			Flash: "",
		})
	}
}

// VerifyUser authenticates a user
func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Authenticate the user
	user := &models.User{}
	user, err = app.DB.AuthenticateUser(r.PostForm.Get("username"), r.PostForm.Get("password"))

	if user == (&models.User{}) {

		session.AddFlash("Invalid Login", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Redirect to login page
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	// Save user info
	session.Values["user"] = user

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Send logged in user to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutUser removes a users session
func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Set user session to empty user
	session.Values["user"] = &models.User{}

	// Save session
	err := session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Send to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// ShowUser display a users info page
func (app *App) ShowUser(w http.ResponseWriter, r *http.Request) {

	// Get requested user
	vars := mux.Vars(r)
	username := vars["username"]

	// Get user from db
	user := &models.User{}
	user, err := app.DB.GetUser(username)
	if err != nil {
		app.NotFound(w)
		return
	}

	// Check for nil user
	if user.ID == 0 {
		app.NotFound(w)
		return
	}

	// Get Reviews
	reviews, err := app.DB.UserLatestReviews(username, 50)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get Collection
	collection, err := app.DB.GetCollection(username)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Get current user
	_, currentUser := app.LoggedIn(r)

	// If a user is viewing their own page
	if username == currentUser.Username {

		// Get current invites from db
		invites, err := app.DB.GetInvites(username)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Display user page with invites
		app.RenderHTML(w, r, "showuser.page.html", &HTMLData{
			DisplayUser: user,
			Invites:     invites,
			Reviews:     reviews,
			Books:       collection,
		})
	} else {

		// Display user page without invites
		app.RenderHTML(w, r, "showuser.page.html", &HTMLData{
			DisplayUser: user,
			Reviews:     reviews,
			Books:       collection,
		})
	}
}

// CreateInviteCode generate an invite code for new users
func (app *App) CreateInviteCode(w http.ResponseWriter, r *http.Request) {

	// Get user info
	user := &models.User{}
	_, user = app.LoggedIn(r)

	// Generate invite code
	code, err := CreateUUID()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Create a new invite
	err = app.DB.CreateInvite(user.Username, code)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Show user their page to display the new invite code
	http.Redirect(w, r, fmt.Sprintf("/user/%s", user.Username), http.StatusSeeOther)
}
