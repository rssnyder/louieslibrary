package models

import (
	"log"
	"database/sql"
)

// InsertMessage
// Send a new message
func (db *DB) InsertMessage(sender, reciver, content string) error {

	// Query statement
	stmt := `INSERT INTO messages (sender, reciver, read, content, created) VALUES ($1, $2, FALSE, $3, timezone('utc', now()))`

	// Create
	row := db.QueryRow(stmt, sender, reciver, content)
	if row != nil {
		err := row.Scan()
		if err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	log.Printf("%s sent a message to %s", sender, reciver)

	// Return id of new review
	return nil
}

// GetConversation
// Retrives messages from a particular user
func (db *DB) GetConversation(sender, reciver string) (Messages, error) {

	// Query statement
	stmt := `SELECT sender, reciver, read, content, created FROM messages WHERE sender = $1 AND reciver = $2 OR  sender = $2 AND reciver = $1 ORDER BY created ASC`

	// Create
	rows, err := db.Query(stmt, sender, reciver)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty message collection
	messages := Messages{}

	// Get all the matching requets
	for rows.Next() {
		m := &Message{}

		// Pull data into message
		err := rows.Scan(&m.Sender, &m.Reciver, &m.Read, &m.Content, &m.Created)
		if err != nil {
			return nil, err
		}

		// Add book to collection
		messages = append(messages, m)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	db.MarkAsRead(reciver, sender)

	// Return id of new review
	return messages, nil
}

func (db *DB) MarkAsRead(sender, reciver string) {

	stmt := `UPDATE messages SET read = TRUE FROM (
		SELECT id FROM  messages WHERE sender = $1 AND reciver = $2) AS subquery
		WHERE messages.id = subquery.id;`

	// Create
	db.QueryRow(stmt, sender, reciver)
}

// GetThreads
// Get the users someone has messages w
func (db *DB) GetThreads(reciver string) (Messages, error) {

	// Unique threads
	var threads []string

	stmt := `SELECT DISTINCT sender FROM messages WHERE reciver = $1`

	// Create
	rows, err := db.Query(stmt, reciver)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty message collection
	messages := Messages{}

	// Get all the matching requets
	for rows.Next() {
		m := &Message{}

		// Pull data into message
		err := rows.Scan(&m.Sender)
		if err != nil {
			return nil, err
		}

		// Add thread if unique
		threads = AppendIfUnique(threads, m.Sender)
	}

	stmt = `SELECT DISTINCT reciver FROM messages WHERE sender = $1`

	// Create
	rows_second, err := db.Query(stmt, reciver)
	if err != nil {
		return nil, err
	}
	defer rows_second.Close()

	// Get all the matching requets
	for rows_second.Next() {
		m := &Message{}

		// Pull data into message
		err := rows_second.Scan(&m.Sender)
		if err != nil {
			return nil, err
		}

		// Add thread if unique
		threads = AppendIfUnique(threads, m.Sender)
	}

	// Catch sql errors
	if err = rows_second.Err(); err != nil {
		return nil, err
	}

	for _, ele := range threads {
		m := &Message{
			Sender: ele,
		}

		// Add message to collection
		messages = append(messages, m)
	}

	// Return id of new review
	return messages, nil
}

// GetUnread
// Get the users unread messages
func (db *DB) GetUnopened(reciver string) (Messages, error) {

	stmt := `SELECT DISTINCT sender FROM messages WHERE reciver = $1 AND read = false`

	// Create
	rows, err := db.Query(stmt, reciver)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Empty message collection
	messages := Messages{}

	// Get all the matching requets
	for rows.Next() {
		m := &Message{}

		// Pull data into message
		err := rows.Scan(&m.Sender)
		if err != nil {
			return nil, err
		}

		// Add book to collection
		messages = append(messages, m)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return id of new review
	return messages, nil
}

// AppendIfUnique
// Append only if item is unique
func AppendIfUnique(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}