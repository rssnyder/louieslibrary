package main

import (
	"net/http"
	"fmt"
	"strconv"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/forms"
)

// ShowBook displays a single book
func (app *App) ShowBook(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Get requested book id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	// Get book
	book, err := app.DB.GetBook(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if book == nil {
		app.NotFound(w)
		return
	}

	// Get Reviews
	reviews, err := app.DB.LatestReviews(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	

	// Get the previous flashes, if any.
	if flashes := session.Flashes("default"); len(flashes) > 0 {
		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.RenderHTML(w, r, "showbook.page.html", &HTMLData{
			Book:   book,
			Form:   &forms.NewReview{},
			Reviews: reviews,
			Flash:  fmt.Sprintf("%v", flashes[0]),
		})
	} else {
		app.RenderHTML(w, r, "showbook.page.html", &HTMLData{
			Book:   book,
			Reviews: reviews,
			Form:   &forms.NewReview{},
			Flash:  "",
		})
	}
}


// NewBook displays the new book upload form
func (app *App) NewBook(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newbook.page.html", &HTMLData{
		Form: &forms.NewBook{},
	})
}

// CreateBook uses a form to create a new book
func (app *App) CreateBook(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	//Limit upload to 10mb
	r.ParseMultipartForm(10 << 20)

	_, user := app.LoggedIn(r)

	// Model the new user book on html form
	form := &forms.NewBook{
		ISBN:         r.PostForm.Get("isbn"),
		Author:       r.PostForm.Get("author"),
		Uploader:      user.Username,
		Title:        r.PostForm.Get("title"),
		Description:  r.PostForm.Get("description"),
		Genre:        r.PostForm.Get("genre"),
	}

	// Validate form
	if !form.Valid() {
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}

	// FormFile returns the first file for the given key `myFile`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
    file, handler, err := r.FormFile("epub")
    if err != nil {
		log.Printf("File Upload Error - %s\n", err)
		form.Upload = false
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
        return
    }
	defer file.Close()
	
    log.Printf("Uploaded File - %+v\n", handler.Filename)
    log.Printf("File Size - %+v\n", handler.Size)
    log.Printf("MIME Header - %+v\n", handler.Header)

    // read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
		log.Printf("File Upload Error - %s\n", err)
		form.Upload = false
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
        return
    }
    // write this byte array to our temporary file
    ioutil.WriteFile(fmt.Sprintf("assets/books/%s - %s.mobi", form.Author, form.Title), fileBytes, 0777)

	// Insert the new book
	id, err := app.DB.InsertBook(form.ISBN, form.Author, form.Uploader, form.Title, form.Description, form.Genre)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Your book was added successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)
}

// CreateReview submits a review for a book
func (app *App) CreateReview(w http.ResponseWriter, r *http.Request) {
	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Model the new review on html form
	form := &forms.NewReview{
		BookID:    r.PostForm.Get("bookid"),
		Username:  r.PostForm.Get("username"),
		Rating:    r.PostForm.Get("rating"),
		Review:    r.PostForm.Get("review"),
	}

	// Validate form
	if !form.Valid() {
		// Save success message
		session.AddFlash("Unable to submit your review.", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/book/%s", r.PostForm.Get("bookid")), http.StatusSeeOther)
		return
	}

	// Insert the new review
	_, err = app.DB.InsertReview(form.BookID, form.Username, form.Rating, form.Review)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Save success message
	session.AddFlash("Your review was added successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/book/%s", form.BookID), http.StatusSeeOther)
}