package models

import (
	"database/sql"
	"log"

	"github.com/rssnyder/louieslibrary/pkg/forms"
)

// GetAnnouncement gets the latest from the db
func (db *DB) GetAnnouncement() (*Announcement, error) {
	// Query statement
	stmt := `SELECT author, content, created FROM announcements WHERE active = TRUE ORDER BY created DESC LIMIT 1`

	// Execute query
	row := db.QueryRow(stmt)
	a := &Announcement{}

	// Pull data into request
	err := row.Scan(&a.Author, &a.Content, &a.Created)
	if err == sql.ErrNoRows {
		log.Printf("Nothing to return")
		return a, nil
	} else if err != nil {
		return a, err
	}

	log.Printf("found something")
	return a, nil
}

// InsertAnnouncement creates a new announcement
func (db *DB) InsertAnnouncement(newAnnouncement *forms.NewAnnouncement) (int, error) {

	// Save stored request
	var id int

	// Query statement
	stmt := `INSERT INTO announcements (author, content, active, created) 
		VALUES ($1, $2, FALSE, timezone('utc', now())) RETURNING id`

	// Query and fill book structure
	err := db.QueryRow(stmt, newAnnouncement.Author, newAnnouncement.Content).Scan(&id)
	if err != nil {
		return 0, err
	}

	log.Printf("New announcement %d uploaded by %s", id, newAnnouncement.Author)

	// Return new announcement id
	return id, nil
}

// RemoveAnnouncement sets an annoucement to not display 
func (db *DB) RemoveAnnouncement() {

	// Query statement
	stmt := `UPDATE announcements SET active = FALSE WHERE active = TRUE`

	log.Printf("Clearing announcements cleared")

	// Update book
	db.QueryRow(stmt)
}
