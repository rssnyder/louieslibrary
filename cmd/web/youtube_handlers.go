package main

import (
	"net/http"
	"fmt"
	"os/exec"
	"os"
)

// NewPlaylist displays the new playlist form
func (app *App) NewPlaylist(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newplaylist.page.html", &HTMLData{})
}

// DownloadPlaylist uses youtubedl to download from yt
func (app *App) DownloadPlaylist(w http.ResponseWriter, r *http.Request) {
	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	r.PostForm.Get("playlisturl")

	uuid, err := CreateUUID()
	if err != nil {
		app.ServerError(w, err)
	}

	// Create directory for output
	savedir := fmt.Sprintf("%s/%s", app.YoutubeDir, uuid)
	os.MkdirAll(savedir, 0777)

	audio_format := "mp3"
	output_format := savedir + "/%(title)s.%(ext)s"	
	_, err = exec.Command("youtube-dl", "--extract-audio", "--audio-format", audio_format, "-i", "-o", output_format, r.PostForm.Get("playlisturl")).Output()
	if err != nil {
		app.ServerError(w, err)
	}

	_, err = ZipDirectory(savedir)
	if err != nil {
		app.ServerError(w, err)
	}
		
	http.Redirect(w, r, fmt.Sprintf("/youtube/%s.zip", uuid), http.StatusSeeOther)
}
