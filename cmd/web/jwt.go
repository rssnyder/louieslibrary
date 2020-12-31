package main

import (
	// "archive/zip"
	// "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	// "fmt"
	// "io/ioutil"
	"net/http"
	// "os"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-cmp/cmp"
	"github.com/rssnyder/louieslibrary/pkg/models"
)

type TokenInfo struct {
	Valid    bool  `json:"valid"`
	TimeLeft int64 `json:"time_left"`
}

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// UserToken is the return structure for a requested jwt
type UserToken struct {
	Token string `json:"token"`
}

// SignJWT Return a signed JWT for a user
func (app *App) SignJWT(username, role string) (UserToken, error) {

	var returnToken UserToken

	// Create claims for user
	claims := CustomClaims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Unix() + 86400,
			Issuer:    "louieslibrary",
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	signedToken, err := token.SignedString(app.SecureString)
	if err != nil {
		return returnToken, err
	}

	returnToken.Token = signedToken

	return returnToken, nil
}

// VerifyJWT Verify a JWT of a user
func (app *App) VerifyJWT(userJwt string) (string, string, int64, error) {

	var username, role string

	// Get claims from token
	token, err := jwt.ParseWithClaims(
		userJwt,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return app.SecureString, nil
		},
	)
	if err != nil {
		return username, role, 0, err
	}

	// Parse the claims from the token
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return username, role, 0, errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return username, role, 0, errors.New("JWT is expired")
	}

	username = claims.Username
	role = claims.Role

	return username, role, (claims.ExpiresAt - time.Now().UTC().Unix()), nil
}

// GetTokenHeader get jwt from request
func GetTokenHeader(r *http.Request) string {

	// Data to return
	var tokenData string

	// Grab authentication header
	token := r.Header["Authorization"]

	// Check for authorization header
	if len(token) > 0 {

		// Check for correct format
		splits := strings.Split(token[0], " ")
		if len(splits) > 0 {

			return splits[1]
		}
	}

	return tokenData
}

// ValidateToken validates if a request has a valid token
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

// ValidateRequest validates if a request has a valid token
func (app *App) ValidateRequest(w http.ResponseWriter, r *http.Request) bool {

	// Get token from header
	token := GetTokenHeader(r)
	if token == "" {
		return false
	}

	// Verify valitity of token
	_, _, _, err := app.VerifyJWT(token)
	if err != nil {
		return false
	}

	// Token valid, return time left
	return true
}

// GetJWT authenticates a user
func (app *App) GetJWT(w http.ResponseWriter, r *http.Request) {

	// Get basic auth
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	// Authenticate the user
	user := &models.User{}
	fail := &models.User{}
	user, err := app.DB.AuthenticateUser(pair[0], pair[1])

	if cmp.Equal(user, fail) {

		// Invalid login attempt
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("{'error':'yes'}")
		return
	}

	// Get signed JWT
	token, err := app.SignJWT(user.Username, user.Role)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	JSONResponse(w, 200, token)
}
