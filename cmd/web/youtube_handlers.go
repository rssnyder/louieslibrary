package main

import (
	"net/http"
	"fmt"
	"os/exec"
	"crypto/rand"
	"os"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/forms"
)

// NewPlaylist displays the new playlist form
func (app *App) NewPlaylist(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newplaylist.page.html", &HTMLData{
		Form: &forms.NewBook{},
	})
}

// DownloadPlaylist uses youtubedl to download from yt
func (app *App) DownloadPlaylist(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	r.PostForm.Get("playlisturl")

	// Create UUId
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
			b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)

	// Create directory for output
	savedir := fmt.Sprintf("%s/%s", app.YoutubeDir, uuid)
	os.MkdirAll(savedir, 0777)

	audio_format := "mp3"
	output_format := savedir + "/%(title)s.%(ext)s"
	fmt.Printf(output_format)
	
	out, err := exec.Command("youtube-dl", "--extract-audio", "--audio-format", audio_format, "-i", "-o", output_format, r.PostForm.Get("playlisturl")).Output()
	if err != nil {
			fmt.Printf("%s", err)
	}
	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)

	output = ZipDirectory(savedir)
	if output == "" {
		// Save message
		session.AddFlash("Unable to download playlist.", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
		
	http.Redirect(w, r, fmt.Sprintf("/youtube/%s.zip", uuid), http.StatusSeeOther)
}
