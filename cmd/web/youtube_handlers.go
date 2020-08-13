package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// NewPlaylist display the playlist form
func (app *App) NewPlaylist(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newplaylist.page.html", &HTMLData{})
}

// DownloadPlaylist use youtubedl to download from yt
func (app *App) DownloadPlaylist(w http.ResponseWriter, r *http.Request) {

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	r.PostForm.Get("playlisturl")

	// Create guid for saving playlist
	uuid, err := CreateUUID()
	if err != nil {
		app.ServerError(w, err)
	}

	// Create directory for output
	savedir := fmt.Sprintf("%s/%s", app.YoutubeDir, uuid)
	os.MkdirAll(savedir, 0777)

	// Use youtube-dl to get playlist in mp3 format
	audioFormat := "mp3"
	outputFormat := savedir + "/%(title)s.%(ext)s"
	_, err = exec.Command("youtube-dl", "--extract-audio", "--audio-format", audioFormat, "-i", "-o", outputFormat, r.PostForm.Get("playlisturl")).Output()
	if err != nil {
		app.ServerError(w, err)
	}

	// Zip the playlist
	fullPath, err := ZipDirectory(savedir)
	if err != nil {
		app.ServerError(w, err)
	}

	// Save playlist to s3 for archival
	err = app.UploadFile("youtube", fmt.Sprintf("playlists/%s.zip", uuid), fullPath)
	if err != nil {
		log.Printf("Unable to send playlist zip to storage %s", err.Error())
	}

	// Send user to playlist file
	http.Redirect(w, r, fmt.Sprintf("/youtube/%s.zip", uuid), http.StatusSeeOther)
}
