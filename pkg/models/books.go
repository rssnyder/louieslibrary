package models

import (
	"database/sql"
	"log"
)

// BooksDB holds the db connection
type BooksDB struct {
	*sql.DB
}

// GetBook retrives a book from the db
func (db *DB) GetBook(id int) (*Book, error) {
	// Query statement
	stmt := `SELECT id, isbn, title, author, uploader, description, genre, created FROM books WHERE id = $1`

	// Execute query
	row := db.QueryRow(stmt, id)
	b := &Book{}

	// Pull data into request
	err := row.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Uploader, &b.Description, &b.Genre, &b.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

// LatestBooks grabs the latest 10 valid books
func (db *DB) LatestBooks(limit int) (Books, error) {
	// Query statement
	stmt := `SELECT id, isbn, title, author, description, genre, created FROM books ORDER BY created DESC LIMIT $1`

	// Execute query
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := Books{}

	// Get all the matching requets
	for rows.Next() {
		b := &Book{}

		// Pull data into request
		err := rows.Scan(&b.ID, &b.ISBN, &b.Title, &b.Author, &b.Description, &b.Genre, &b.Created)
		if err != nil {
			return nil, err
		}

		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// InsertBooks adds a new book to the library
func (db *DB) InsertBook(isbn, title, author, uploader, description, genre string) (int, error) {
	// Save stored request
	var bookid int

	// Query statement
	stmt := `INSERT INTO books (isbn, title, author, uploader, description, genre, created) VALUES ($1, $2, $3, $4, $5, $6, timezone('utc', now())) RETURNING id`

	err := db.QueryRow(stmt, isbn, title, author, uploader, description, genre).Scan(&bookid)
	if err != nil {
		return 0, err
	}

	log.Printf("New book %s uploaded by %s", title, uploader)

	return bookid, nil
}
