package models

import (
	"database/sql"
	"log"

	"github.com/Mr-Schneider/louieslibrary/pkg/forms"
)

// GetBook retrive a book from the db
func (db *DB) GetBook(id string) (*Book, error) {
	// Query statement
	stmt := `SELECT id, volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created FROM books WHERE volumeid = $1`

	// Execute query
	row := db.QueryRow(stmt, id)
	b := &Book{}

	// Pull data into request
	err := row.Scan(&b.ID, &b.VolumeID, &b.Title, &b.Subtitle, &b.Publisher, &b.PublishedDate, &b.PageCount,
		&b.MaturityRating, &b.Authors, &b.Categories, &b.Description, &b.Uploader, &b.Price, &b.ISBN10, &b.ISBN13,
		&b.ImageLink, &b.Downloads, &b.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

// LatestBooks grab the latest n books
func (db *DB) LatestBooks(limit int) (Books, error) {
	// Query statement
	stmt := `SELECT id, volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created FROM books ORDER BY created DESC LIMIT $1`

	// Execute query
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty book collection
	books := Books{}

	// Get all the matching requets
	for rows.Next() {
		b := &Book{}

		// Pull data into request
		err := rows.Scan(&b.ID, &b.VolumeID, &b.Title, &b.Subtitle, &b.Publisher, &b.PublishedDate, &b.PageCount,
			&b.MaturityRating, &b.Authors, &b.Categories, &b.Description, &b.Uploader, &b.Price, &b.ISBN10, &b.ISBN13,
			&b.ImageLink, &b.Downloads, &b.Created)
		if err != nil {
			return nil, err
		}

		// Add book to collection
		books = append(books, b)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return collection of recent books
	return books, nil
}

// InsertBook add a new book to the library
func (db *DB) InsertBook(newBook *forms.NewBook) (int, error) {

	// Save stored request
	var bookid int

	// Query statement
	stmt := `INSERT INTO books (volumeid, title, subtitle, publisher, publisheddate, pagecount,
		maturityrating, authors, categories, description, uploader, price, isbn10, isbn13,
		imagelink, downloads, created) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 0, timezone('utc', now())) RETURNING id`

	// Query and fill book structure
	err := db.QueryRow(stmt, newBook.VolumeID, newBook.Title, newBook.Subtitle, newBook.Publisher, newBook.PublishedDate, newBook.PageCount,
		newBook.MaturityRating, newBook.Authors, newBook.Categories, newBook.Description, newBook.Uploader, newBook.Price, newBook.ISBN10, newBook.ISBN13,
		newBook.ImageLink).Scan(&bookid)
	if err != nil {
		return 0, err
	}

	log.Printf("New book %s uploaded by %s", newBook.Title, newBook.Uploader)

	// Return new book id
	return bookid, nil
}

// DownloadBook increment downloads of a book
func (db *DB) DownloadBook(bookID string, downloads int) {

	// Query statement
	stmt := `UPDATE books SET downloads = $1 WHERE volumeid = $2`

	// Incriment count
	db.QueryRow(stmt, downloads, bookID)
}

// UpdateBook edit a books attributes
func (db *DB) UpdateBook(book *forms.NewBook) (int, error) {

	// Save stored request
	var bookid int

	// Query statement
	stmt := `UPDATE books SET title = $1, subtitle = $2, publisher = $3, publisheddate = $4, pagecount = $5,
		maturityrating = $6, authors = $7, categories = $8, description = $9, price = $10, isbn10 = $11, isbn13 = $12, imagelink = $13 WHERE volumeid = $14`

	// Update book
	db.QueryRow(stmt, book.Title, book.Subtitle, book.Publisher, book.PublishedDate, book.PageCount,
		book.MaturityRating, book.Authors, book.Categories, book.Description, book.Price, book.ISBN10, book.ISBN13,
		book.ImageLink, book.VolumeID)

	log.Printf("Book %s edited", book.Title)

	// Return book id of edited book
	return bookid, nil
}

// CollectBook add a book to a users collection
func (db *DB) CollectBook(username, year, id string) {

	// Query statement
	stmt := `INSERT INTO collection (username, volumeid, year, created) VALUES ($1, $3, $2, timezone('utc', now()))`

	db.QueryRow(stmt, username, year, id)

	log.Printf("%s collected book %s", username, id)
}

// GetCollectionItem get a book from a users collection
func (db *DB) GetCollectionItem(username, id string) bool {

	var vol int

	// Query statement
	stmt := `SELECT volumeid FROM collection WHERE username = $1 AND volumeid = $2`

	err := db.QueryRow(stmt, username, id).Scan(&vol)
	if err == sql.ErrNoRows {
		return false
	}

	return true
}

// GetCollection get all the books in a users collection
func (db *DB) GetCollection(username string) (Books, error) {

	// Query statement
	stmt := `SELECT c.volumeid id, b.title, b.imagelink, c.year FROM collection c 
		INNER JOIN books b ON c.volumeid = b.volumeid AND c.username = $1 ORDER BY c.year DESC`

	// Execute query
	rows, err := db.Query(stmt, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty book collection
	books := Books{}

	// Get all the matching requets
	for rows.Next() {
		b := &Book{}

		// Pull data into request
		err := rows.Scan(&b.VolumeID, &b.Title, &b.ImageLink, &b.Subtitle)
		if err != nil {
			return nil, err
		}

		// Add book to collection
		books = append(books, b)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Got collection for %s", username)

	// Return collection for user
	return books, nil
}
