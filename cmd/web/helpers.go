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

func ZipDirectory(dirPath string) (string, error) {
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

	err = writer.Close()
	if err != nil {
		return output, err
	}

	return dirPath + ".zip", nil
}

func CreateUUID() (string, error) {
	var uuid string

	// Create UUId
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return uuid, errors.New("Unable to generate UUID")
	}

	uuid = fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}