package main

import (
	"fmt"
	"errors"
	"net/http"
	"io/ioutil"
	"archive/zip"
	"os"
	"crypto/rand"
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
)

// LoggedIn
// Get logged in status
// Pass along user data if true
func (app *App) LoggedIn(r *http.Request) (bool, *models.User) {

	// Empty user struct
	var user = &models.User{}

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get user from session
	logged_in := session.Values["user"]
	user, ok := logged_in.(*models.User)

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

// ZipDirectory
// Compress a directory on the disk
// Return the name of the zip file
func ZipDirectory(dirPath string) (string, error) {

	// Empty output
	var output string

	// Get a Buffer to Write To
	out_file, err := os.Create(dirPath + ".zip")
	if err != nil {
			return output, err
	}
	defer out_file.Close()

	// Create a new zip archive.
	writer := zip.NewWriter(out_file)

	// Add all file sin directory to zip
	files, err := ioutil.ReadDir(dirPath)
	for _, file := range files {
		full_path := dirPath + "/" + file.Name()
		dat, err := ioutil.ReadFile(full_path)
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

// CreateUUID
// Generate a guid
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

// AppendIfUnique
// Append only if item is unique
func AppendIfUnique(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}