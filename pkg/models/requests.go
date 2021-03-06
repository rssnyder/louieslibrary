package models

import (
	"database/sql"
	"log"
)

// GetRequest retrive a request from the db
func (db *DB) GetRequest(id int) (*Request, error) {

	// Query statement
	stmt := `SELECT id, requester, title, status, bookid, created FROM requests WHERE id = $1`

	// Execute query
	row := db.QueryRow(stmt, id)
	r := &Request{}

	// Pull data into request
	err := row.Scan(&r.ID, &r.Requester, &r.Title, &r.Status, &r.BookID, &r.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	log.Printf("BOOKID IS %s", r.BookID)

	// Return review
	return r, nil
}

// LatestRequests grab latest n requests
func (db *DB) LatestRequests(limit int) (Requests, error) {

	// Query statement
	stmt := `SELECT id, requester, title, status, created FROM requests ORDER BY created DESC LIMIT $1`

	// Execute query
	rows, err := db.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty request collection
	requests := Requests{}

	// Get all the matching requets
	for rows.Next() {
		r := &Request{}

		// Pull data into request
		err := rows.Scan(&r.ID, &r.Requester, &r.Title, &r.Status, &r.Created)
		if err != nil {
			return nil, err
		}

		// Add request to collection
		requests = append(requests, r)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

// InsertRequest add new request to the db
func (db *DB) InsertRequest(requester, title, source string) (int, error) {

	// Save stored request
	var requestid int

	// Query statement
	stmt := `INSERT INTO requests (requester, title, source, status, bookid, created) VALUES ($1, $2, $3, 'missing', '', timezone('utc', now())) RETURNING id`

	// Create new request
	err := db.QueryRow(stmt, requester, title, source).Scan(&requestid)
	if err != nil {
		return 0, err
	}

	log.Printf("New request submitted by %s", requester)

	// Return new request id
	return requestid, nil
}

// FillRequest link a request to a book
func (db *DB) FillRequest(requestid int, bookid string) string {

	// Query statement
	stmt := `UPDATE requests SET bookid = $1, status = 'found' WHERE id = $2`

	// Link
	err := db.QueryRow(stmt, bookid, requestid)
	if err == nil {
		return ""
	}

	log.Printf("Request %d filled", requestid)

	return bookid
}
