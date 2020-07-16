package main

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims holds the template for a jwt
type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// UserToken is the return structure for a requested jwt
type UserToken struct {
	Token string `json:"token"`
}

// LoggedIn Get logged in status
// Pass along user data if true
func (app *App) LoggedIn(r *http.Request) (bool, *models.User) {

	// Empty user struct
	var user = &models.User{}

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get user from session
	loggedIn := session.Values["user"]
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

// ZipDirectory Compress a directory on the disk
// Return the name of the zip file
func ZipDirectory(dirPath string) (string, error) {

	// Empty output
	var output string

	// Get a Buffer to Write To
	outFile, err := os.Create(dirPath + ".zip")
	if err != nil {
		return output, err
	}
	defer outFile.Close()

	// Create a new zip archive.
	writer := zip.NewWriter(outFile)

	// Add all file sin directory to zip
	files, err := ioutil.ReadDir(dirPath)
	for _, file := range files {
		fullPath := dirPath + "/" + file.Name()
		dat, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return output, err
		}

		// Add some files to the archive.
		f, err := writer.Create(file.Name())
		if err != nil {
			return output, err
		}
		_, err = f.Write(dat)
		if err != nil {
			return output, err
		}
	}

	// Attempt to finish the file
	err = writer.Close()
	if err != nil {
		return output, err
	}

	// Return the zip file location
	return dirPath + ".zip", nil
}

// CreateUUID Generate a guid
func CreateUUID() (string, error) {

	// Empty guid
	var uuid string

	// Create UUID
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return uuid, errors.New("Unable to generate UUID")
	}

	// Format giud
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	// Return the guid
	return uuid, nil
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
