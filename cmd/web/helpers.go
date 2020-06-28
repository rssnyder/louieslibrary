package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"archive/zip"
	"os"

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

func ZipDirectory(dirPath string) string {
	// Get a Buffer to Write To
	outFile, err := os.Create(dirPath + ".zip")
	if err != nil {
			fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	writer := zip.NewWriter(outFile)

	files, err := ioutil.ReadDir(dirPath)
	for _, file := range files {
		fullPath := dirPath + "/" + file.Name()
		log.Printf("reading file %s to zip", fullPath)
		dat, err := ioutil.ReadFile(fullPath)
		if err != nil {
			log.Printf("failed reading file %s to zip, %s", fullPath, err)
			return ""
		}

		// Add some files to the archive.
		f, err := writer.Create(file.Name())
		if err != nil {
			log.Printf("failed creating file %s in zip, %s", fullPath, err)
			return ""
		}
		_, err = f.Write(dat)
		if err != nil {
			log.Printf("failed writing file %s in zip, %s", fullPath, err)
			return ""
		}
	}

	err = writer.Close()
	if err != nil {
		log.Printf("failed closing zip writer %s, %s", dirPath, err)
		return ""
	}

	return dirPath + ".zip"
}