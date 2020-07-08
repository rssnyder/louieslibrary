package main

import (
	"net/http"
	"fmt"
	"strconv"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"
	"path/filepath"
	"github.com/Mr-Schneider/louieslibrary/pkg/forms"
)

// ShowBook
// Display a page with information on a single book
func (app *App) ShowBook(w http.ResponseWriter, r *http.Request) {

	// Get requested book id
	vars := mux.Vars(r)
	id := vars["volumeid"]

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
	reviews, err := app.DB.LatestReviews(id, 50)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Render page without flash
	app.RenderHTML(w, r, "showbook.page.html", &HTMLData{
		Book:   book,
		Reviews: reviews,
		Form:   &forms.NewReview{},
	})
}

// DownloadBook
// Send the user the file for the requested book
func (app *App) DownloadBook(w http.ResponseWriter, r *http.Request) {

	// Get requested book id
	vars := mux.Vars(r)
	id := vars["volumeid"]

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

	// Perform nessesary actions on book being downloaded
	app.DB.DownloadBook(book.VolumeID, book.Downloads + 1)

	// Present book file to user
	app.ServeFile(w, app.BookBucket, fmt.Sprintf("%s.mobi", book.VolumeID), fmt.Sprintf("%s - %s.mobi", book.Title, book.Authors))
}

// NewBook
// Display the new book form
func (app *App) NewBook(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "newbook.page.html", &HTMLData{
		Form: &forms.NewBook{},
	})
}

// CreateBook
// Build the new book structure
// If only volumeid given, grab data from books.google api
// Otherwise user given information
func (app *App) CreateBook(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	//Limit upload to 50mb
	r.ParseMultipartForm(50 << 20)

	// Get current user
	_, user := app.LoggedIn(r)

	// Grab information from books.google if no title given
	if r.PostForm.Get("title") == "" {

		// Grab volumeid from form
		book_info := GetBookInfo(r.PostForm.Get("volumeid"), app.BookAPIKey)

		// Model the new book on api feedback
		form := &forms.NewBook{
			VolumeID:				book_info.Id,
			Title:       		book_info.Data.Title,
			Subtitle:				book_info.Data.Subtitle,
			Publisher:			book_info.Data.Publisher,
			PublishedDate:	book_info.Data.PublishedDate,
			PageCount:			strconv.Itoa(book_info.Data.PageCount),
			MaturityRating:	book_info.Data.MaturityRating,
			Authors:      	fmt.Sprint(book_info.Data.Authors),
			Categories:     fmt.Sprint(book_info.Data.Categories),
			Description: 		book_info.Data.Description,
			Uploader: 			user.Username,
			Price:					fmt.Sprintf("%.2f %s", book_info.SaleInfo.Retail.Amount, book_info.SaleInfo.Retail.CurrencyCode),
			ISBN10:					fmt.Sprintf("%s %s", book_info.Data.IndustryIdentifiers[0].Type, book_info.Data.IndustryIdentifiers[0].Identifier),
			ISBN13:					fmt.Sprintf("%s %s", book_info.Data.IndustryIdentifiers[0].Type, book_info.Data.IndustryIdentifiers[1].Identifier),
			ImageLink:			fmt.Sprint(book_info.Data.ImageLinks.Small),
		}

		// Display the new book form with the retrived data
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}

	// Model the new book on the information from the form
	form := &forms.NewBook{
		VolumeID:				r.PostForm.Get("volumeid"),
		Title:       		r.PostForm.Get("title"),
		Subtitle:				r.PostForm.Get("subtitle"),
		Publisher:			r.PostForm.Get("publisher"),
		PublishedDate:	r.PostForm.Get("publisheddate"),
		PageCount:			r.PostForm.Get("pagecount"),
		MaturityRating:	r.PostForm.Get("maturityrating"),
		Authors:      	r.PostForm.Get("authors"),
		Categories:     r.PostForm.Get("categories"),
		Description: 		r.PostForm.Get("description"),
		Uploader: 			user.Username,
		Price:					r.PostForm.Get("price"),
		ISBN10:					r.PostForm.Get("isbn10"),
		ISBN13:					r.PostForm.Get("isbn13"),
		ImageLink:			r.PostForm.Get("imagelink"),
	}

	// Validate the new book form
	if !form.Valid() {
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}

	// Get book file from form
	file, handler, err := r.FormFile("epub")
	if err != nil {
		log.Printf("File Upload Error - %s\n", err)
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}
	defer file.Close()

	// Read contents of uploaded file
	file_bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("File Upload Error - %s\n", err)
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}

	// Get uploaded file format
	format := filepath.Ext(handler.Filename)

	// Send book to storage server
	err = app.UploadBytes(app.BookBucket, fmt.Sprintf("%s.%s", form.VolumeID, format), file_bytes)
	if err != nil {
		app.ServerError(w, err)
	}

	// Insert the new book
	_, err = app.DB.InsertBook(form)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session.AddFlash("Your book was added successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Direct to new book page
	http.Redirect(w, r, fmt.Sprintf("/book/%s", form.VolumeID), http.StatusSeeOther)
}

// CreateReview
// Build the new review structure and submit
func (app *App) CreateReview(w http.ResponseWriter, r *http.Request) {

	// Load session
	session, _ := app.Sessions.Get(r, "session-name")

	// Parse the post data
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Get the current user
	_, user := app.LoggedIn(r)

	// Model the new review on html form
	form := &forms.NewReview{
		BookID:    r.PostForm.Get("volumeid"),
		Username:  user.Username,
		Rating:    r.PostForm.Get("rating"),
		Review:    r.PostForm.Get("review"),
	}

	// Validate the new review form
	if !form.Valid() {

		session.AddFlash("Unable to submit your review.", "default")

		// Save session
		err = session.Save(r, w)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Send back to book page
		http.Redirect(w, r, fmt.Sprintf("/book/%s", r.PostForm.Get("bookid")), http.StatusSeeOther)
		return
	}

	// Insert the new review
	_, err = app.DB.InsertReview(form.BookID, form.Username, form.Rating, form.Review)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	session.AddFlash("Your review was added successfully!", "default")

	// Save session
	err = session.Save(r, w)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Send back to book page
	http.Redirect(w, r, fmt.Sprintf("/book/%s", form.BookID), http.StatusSeeOther)
}

// ListAllBooks
// Display a list off all books
func (app *App) ListAllBooks(w http.ResponseWriter, r *http.Request) {

	// Get the books
	books, err := app.DB.LatestBooks(1000)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Display page with all books
	app.RenderHTML(w, r, "showbooks.page.html", &HTMLData{
		Books:    books,
	})
}

// EditBook
// Show the new book form
// Filled with the current book data
func (app *App) EditBook(w http.ResponseWriter, r *http.Request) {

	// Get requested book id
	vars := mux.Vars(r)
	id := vars["volumeid"]

	// Get current book data
	book, err := app.DB.GetBook(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if book == nil {
		app.NotFound(w)
		return
	}

	// Model the book
	form := &forms.NewBook{
		ID:							book.ID,
		VolumeID:				book.VolumeID,
		Title:       		book.Title,
		Subtitle:				book.Subtitle,
		Publisher:			book.Publisher,
		PublishedDate:	book.PublishedDate,
		PageCount:			book.PageCount,
		MaturityRating:	book.MaturityRating,
		Authors:      	book.Authors,
		Categories:     book.Categories,
		Description: 		book.Description,
		Price:					book.Price,
		ISBN10:					book.ISBN10,
		ISBN13:					book.ISBN13,
		ImageLink:			book.ImageLink,
	}

	// Display the new book page with the current data
	app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
}

// UpdateBook
// Process a new book form and update an existing book
func (app *App) UpdateBook(w http.ResponseWriter, r *http.Request) {

	//Limit upload to 50mb
	r.ParseMultipartForm(50 << 20)

	form := &forms.NewBook{
		VolumeID:				r.PostForm.Get("volumeid"),
		Title:       		r.PostForm.Get("title"),
		Subtitle:				r.PostForm.Get("subtitle"),
		Publisher:			r.PostForm.Get("publisher"),
		PublishedDate:	r.PostForm.Get("publisheddate"),
		PageCount:			r.PostForm.Get("pagecount"),
		MaturityRating:	r.PostForm.Get("maturityrating"),
		Authors:      	r.PostForm.Get("authors"),
		Categories:     r.PostForm.Get("categories"),
		Description: 		r.PostForm.Get("description"),
		Uploader: 			"no change",
		Price:					r.PostForm.Get("price"),
		ISBN10:					r.PostForm.Get("isbn10"),
		ISBN13:					r.PostForm.Get("isbn13"),
		ImageLink:			r.PostForm.Get("imagelink"),
	}

	// Validate new book form
	if !form.Valid() {
		app.RenderHTML(w, r, "newbook.page.html", &HTMLData{Form: form})
		return
	}
	
	// Update the book with the new information
	app.DB.UpdateBook(form)

	// Display the edited book page
	http.Redirect(w, r, fmt.Sprintf("/book/%s", form.VolumeID), http.StatusSeeOther)
}

// AddToCollection
// Add a book to the current users collection
// Requires the target books volumeid as a var
func (app *App) AddToCollection(w http.ResponseWriter, r *http.Request) {

	// Get requested book id
	vars := mux.Vars(r)
	id := vars["volumeid"]

	// Get current user
	_, user := app.LoggedIn(r)

	// Add book to users collection
	app.DB.CollectBook(user.Username, id)

	// Display the added books page
	http.Redirect(w, r, fmt.Sprintf("/book/%s", id), http.StatusSeeOther)
}