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
func (db *DB) GetUser(username string) (*User, error) {
	// Get attributes of user
	row := db.QueryRow("SELECT id, username, role, created FROM users WHERE username = $1", username)

	u := &User{}

	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Created)
	if err == sql.ErrNoRows {
		return &User{}, nil
	} else if err != nil {
		return &User{}, err
	}

	log.Printf("Retrived data on user %s", username)

	return u, nil
}

// GetInvites checks the valitity of an invite code
func (db *DB) GetInvites(creator string) (Invites, error) {
	// Query statement
	stmt := `SELECT id, code, username, creator, status, created FROM invites WHERE creator = $1 ORDER BY created DESC`

	// Execute query
	rows, err := db.Query(stmt, creator)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invites := Invites{}

	// Get all the matching requets
	for rows.Next() {
		i := &Invite{}

		// Pull data into request
		err := rows.Scan(&i.ID, &i.Code, &i.Username, &i.Creator, &i.Activated, &i.Created)
		if err != nil {
			return nil, err
		}

		invites = append(invites, i)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invites, nil
}

// ValidateInvite checks the valitity of an invite code
func (db *DB) ValidateInvite(invite_code string) (bool, error) {

	// Get id and password hash for given username
	row := db.QueryRow("SELECT status FROM invites WHERE code = $1", invite_code)

	var status bool

	err := row.Scan(&status)
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return true, nil
	}

	return status, nil
}

// CreateInvite creates a new invite
func (db *DB) CreateInvite(creator, code string) error {

	var id int

	stmt := `INSERT INTO invites (code, creator, status, created) VALUES($1, $2, FALSE, timezone('utc', now())) RETURNING id`

	err := db.QueryRow(stmt, code, creator).Scan(&id)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	return nil
}

// FillInvite uses an invite for a new user
func (db *DB) FillInvite(username, code string) error {

	var id int

	stmt := `UPDATE invites SET username = $1, activated = timezone('utc', now()), status = TRUE WHERE code = $2 RETURNING id`

	err := db.QueryRow(stmt, username, code).Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return nil
	}

	return nil
}