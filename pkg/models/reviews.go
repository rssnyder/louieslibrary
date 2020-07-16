package models

import (
	"database/sql"
	"log"
)

// GetReview
// Get a review from the db
func (db *DB) GetReview(bookid string) (*Review, error) {

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

// LatestReviews
// Grab latest n reviews
func (db *DB) LatestReviews(bookid string, limit int) (Reviews, error) {

	// Query statement
	stmt := `SELECT bookid, username, rating, review, created FROM reviews WHERE bookid = $1 ORDER BY created DESC LIMIT $2`

	// Execute query
	rows, err := db.Query(stmt, bookid, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty review collection
	reviews := Reviews{}

	// Get all the matching requets
	for rows.Next() {
		r := &Review{}

		// Pull data into request
		err := rows.Scan(&r.BookID, &r.Username, &r.Rating, &r.Review, &r.Created)
		if err != nil {
			return nil, err
		}

		// Add review to collection
		reviews = append(reviews, r)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// UserLatestReviews
// Get users reviews
func (db *DB) UserLatestReviews(username string, limit int) (Reviews, error) {

	// Query statement
	stmt := `SELECT r.id, r.bookid id, r.rating, r.review, r.created, b.title FROM reviews r 
	INNER JOIN books b ON r.bookid = b.volumeid AND r.username = $1 ORDER BY created DESC LIMIT $2`

	// Execute query
	rows, err := db.Query(stmt, username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty review collection
	reviews := Reviews{}

	// Get all the matching requets
	for rows.Next() {
		r := &Review{}

		// Pull data into request
		err := rows.Scan(&r.ID, &r.BookID, &r.Rating, &r.Review, &r.Created, &r.Username)
		if err != nil {
			return nil, err
		}

		// Add review to collection
		reviews = append(reviews, r)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// InsertReview
// Add a new review
func (db *DB) InsertReview(bookid, username, rating, review string) (int, error) {

	// Save stored request
	var reviewid int

	// Query statement
	stmt := `INSERT INTO reviews (bookid, username, rating, review, created) VALUES ($1, $2, $3, $4, timezone('utc', now())) RETURNING id`

	// Create
	err := db.QueryRow(stmt, bookid, username, rating, review).Scan(&reviewid)
	if err != nil {
		return 0, err
	}

	log.Printf("New review added by %s", username)

	// Return id of new review
	return reviewid, nil
}
