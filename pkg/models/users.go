package models

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// InsertUser creates a new user
func (db *DB) InsertUser(name, email, password string) error {

	var userid int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, password, role, created) VALUES($1, $2, $3, 'reader', timezone('utc', now())) RETURNING id`

	err = db.QueryRow(stmt, name, email, hashedPassword).Scan(&userid)
	if err != nil {
		return err
	}

	log.Printf("User %s registered!", name)

	return nil
}

// AuthenticateUser checks the valitity of a login request
func (db *DB) AuthenticateUser(username, password string) (*User, error) {

	// Get id and password hash for given username
	row := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username)

	u := &User{}

	err := row.Scan(&u.ID, &u.Username, &u.HashedPassword, &u.Role)
	if err == sql.ErrNoRows {
		return &User{}, nil
	} else if err != nil {
		return &User{}, err
	}

	// Check whether the hashed password and plain-text password provided match
	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return &User{}, nil
	} else if err != nil {
		return &User{}, err
	}

	log.Printf("User %s logged in", username)

	return u, nil
}

// GetUser retrives user information
func (db *DB) GetUser(id int) (*User, error) {
	return nil, nil
}
