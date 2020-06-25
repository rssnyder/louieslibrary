package models

import (
	"database/sql"
)

// GetRequest retrives a request from the db
func (db *DB) GetRequest(id int) (*Request, error) {
	// Query statement
	stmt := `SELECT id, requester, title, status, created FROM requests WHERE id = $1`

	// Execute query
	row := db.QueryRow(stmt, id)
	r := &Request{}

	// Pull data into request
	err := row.Scan(&r.ID, &r.Requester, &r.Title, &r.Status, &r.Created)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return r, nil
}

// LatestRequests grabs the latest 10 valid request
func (db *DB) LatestRequests() (Requests, error) {
	// Query statement
	stmt := `SELECT id, requester, title, status, created FROM requests ORDER BY created DESC LIMIT 10`

	// Execute query
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	requests := Requests{}

	// Get all the matching requets
	for rows.Next() {
		r := &Request{}

		// Pull data into request
		err := rows.Scan(&r.ID, &r.Requester, &r.Title, &r.Status, &r.Created)
		if err != nil {
			return nil, err
		}

		requests = append(requests, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

// InsertRequest adds a new request to the db
func (db *DB) InsertRequest(requester, title, source string) (int, error) {
	// Save stored request
	var requestid int

	// Query statement
	stmt := `INSERT INTO requests (requester, title, source, status, created) VALUES ($1, $2, $3, 'missing', timezone('utc', now())) RETURNING id`

	err := db.QueryRow(stmt, requester, title, source).Scan(&requestid)
	if err != nil {
		return 0, err
	}

	return requestid, nil
}
