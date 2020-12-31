package main

import (
	// "archive/zip"
	// "crypto/rand"
	"errors"
	// "fmt"
	// "io/ioutil"
	"net/http"
	// "os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/rssnyder/louieslibrary/pkg/models"
)

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