package main

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rssnyder/louieslibrary/pkg/models"
)

// LoggedIn get logged in status
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

// ZipDirectory compress a directory on the disk
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

// CreateUUID generate a guid
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
