package models

import (
	"database/sql"
)

// GetBook retrives a book from the db
func (db *DB) GetReview(bookid int) (*Review, error) {
	// Query statement
	stmt := `SELECT bookid, username, rating, review, created FROM reviews WHERE bookid = $1`

	// Execute query
	row := db.QueryRow(stmt, bookid)
	r := &Review{}

	// Pull data into request*DB
	err := row.Scan(&r.BookID, &r.Username, &r.Rating, &r.Review, &r.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r, nil
}

// LatestBooks grabs the latest 10 valid books
func (db *DB) LatestReviews(bookid int) (Reviews, error) {
	// Query statement
	stmt := `SELECT bookid, username, rating, review, created FROM reviews WHERE bookid = $1 ORDER BY created DESC LIMIT 10`

	// Execute query
	rows, err := db.Query(stmt, bookid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := Reviews{}

	// Get all the matching requets
	for rows.Next() {
		r := &Review{}

		// Pull data into request
		err := rows.Scan(&r.BookID, &r.Username, &r.Rating, &r.Review, &r.Created)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// InsertBooks adds a new book to the library
func (db *DB) InsertReview(bookid, username, rating, review string) (int, error) {
	// Save stored request
	var reviewid int

	// Query statement
	stmt := `INSERT INTO reviews (bookid, username, rating, review, created) VALUES ($1, $2, $3, $4, timezone('utc', now())) RETURNING id`

	err := db.QueryRow(stmt, bookid, username, rating, review).Scan(&reviewid)
	if err != nil {
		return 0, err
	}

	return reviewid, nil
}
