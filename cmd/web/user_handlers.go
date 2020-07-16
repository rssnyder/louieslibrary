package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Mr-Schneider/louieslibrary/pkg/forms"
	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

// UserLogin holds login data
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUser holds signup data
type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// UserPage holds user data for display
type UserPage struct {
	User    *models.User     `json:"user"`
	Reviews []*models.Review `json:"reviews"`
}

// TokenInfo holds token valitity
type TokenInfo struct {
	Valid    bool  `json:"valid"`
	TimeLeft int64 `json:"time_left"`
}

// CreateUser Use signup form to create a new user
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Store user data
	var userLogin NewUser
	decoder := json.NewDecoder(r.Body)

	// Get login info from request)
	err := decoder.Decode(&userLogin)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Model the new user based on html form
	form := &forms.NewUser{
		Username:   userLogin.Username,
		Email:      userLogin.Email,
		InviteCode: userLogin.Code,
		Password:   userLogin.Password,
	}

	// Validate form
	if !form.Valid() {
		JSONResponse(w, 400, "error")
		return
	}

	// Validate request
	used, err := app.DB.ValidateInvite(form.InviteCode)
	if err != nil {
		JSONResponse(w, 500, "error")
		return
	}
	if used {
		JSONResponse(w, 400, "error")
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

	JSONResponse(w, 200, "success")
}

// VerifyUser Authenticates a user
func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {

	// Store user login
	var userLogin UserLogin
	decoder := json.NewDecoder(r.Body)

	// Get login info from request
	err := decoder.Decode(&userLogin)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Authenticate the user
	user := &models.User{}
	fail := &models.User{}
	user, err = app.DB.AuthenticateUser(userLogin.Username, userLogin.Password)

	if cmp.Equal(user, fail) {

		// Invalid login attempt
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("{'error':'yes'}")
		return
	}

	// Get signed JWT
	token, err := app.SignJWT(user.Username, user.Role)

	JSONResponse(w, 200, token)
}

// ShowUser Display a users info page
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
	// collection, err := app.DB.GetCollection(username)
	// if err != nil {
	// 	app.ServerError(w, err)
	// 	return
	// }

	// // Get current user
	// _, currentUser := app.LoggedIn(r)

	// // If a user is viewing their own page
	// if username == currentUser.Username {

	// 	// Get current invites from db
	// 	invites, err := app.DB.GetInvites(username)
	// 	if err != nil {
	// 		app.ServerError(w, err)
	// 		return
	// 	}

	// 	// Display user page with invites
	// 	app.RenderHTML(w, r, "showuser.page.html", &HTMLData{
	// 		DisplayUser: user,
	// 		Invites:     invites,
	// 		Reviews:     reviews,
	// 		Books:       collection,
	// 	})
	// } else {

	// 	// Display user page without invites
	// 	app.RenderHTML(w, r, "showuser.page.html", &HTMLData{
	// 		DisplayUser: user,
	// 		Reviews:     reviews,
	// 		Books:       collection,
	// 	})
	// }

	userData := UserPage{
		User:    user,
		Reviews: reviews,
	}

	JSONResponse(w, 200, userData)
}

// CreateInviteCode Generate an invite code for new users
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

// ValidateToken returns if a given token is valid
func (app *App) ValidateToken(w http.ResponseWriter, r *http.Request) {

	// Data to return
	var tokenInfo TokenInfo

	// Get token from header
	token := GetTokenHeader(r)
	if token == "" {
		JSONResponse(w, 401, "")
		return
	}

	// Verify valitity of token
	_, _, left, err := app.VerifyJWT(token)
	if err != nil {
		log.Println("Invalid token verify")
		JSONResponse(w, 401, "")
		return
	}

	// Token valid, return time left
	tokenInfo.TimeLeft = left
	tokenInfo.Valid = true
	JSONResponse(w, 200, tokenInfo)
	return
}
