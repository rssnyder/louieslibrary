package models

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser create a new user
func (db *DB) InsertUser(name, email, password string) error {

	// Empty new user id
	var userid int

	// Hash and salt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, password, role, created) VALUES($1, $2, $3, 'reader', timezone('utc', now())) RETURNING id`

	// Create
	err = db.QueryRow(stmt, name, email, hashedPassword).Scan(&userid)
	if err != nil {
		return err
	}

	log.Printf("User %s registered!", name)

	return nil
}

// AuthenticateUser checks the valitity of a login
func (db *DB) AuthenticateUser(username, password string) (*User, error) {

	// Empty user
	u := &User{}

	// Get id and password hash for given username
	row := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username)

	// Pull in password for comparesson
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

	// Return logged in user
	return u, nil
}

// GetUser retrive user information
func (db *DB) GetUser(username string) (*User, error) {

	// Empty user
	u := &User{}

	// Get attributes of user
	row := db.QueryRow("SELECT id, username, role, created FROM users WHERE username = $1", username)

	// Grab user
	err := row.Scan(&u.ID, &u.Username, &u.Role, &u.Created)
	if err == sql.ErrNoRows {
		return &User{}, nil
	} else if err != nil {
		return &User{}, err
	}

	log.Printf("Retrived data on user %s", username)

	return u, nil
}

// GetUsers retrive user information on everyone
func (db *DB) GetUsers() (Users, error) {

	// Empty user
	users := Users{}

	// Get attributes of user
	rows, err := db.Query("SELECT username, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get all the matching requets
	for rows.Next() {
		u := &User{}

		// Pull data into request
		err := rows.Scan(&u.Username, &u.Role)
		if err != nil {
			return nil, err
		}

		// Add invite to collection
		users = append(users, u)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetInvites get a users invites
func (db *DB) GetInvites(creator string) (Invites, error) {

	// Empty invite collection
	invites := Invites{}

	// Query statement
	stmt := `SELECT id, code, username, creator, status, created FROM invites WHERE creator = $1 ORDER BY created DESC`

	// Execute query
	rows, err := db.Query(stmt, creator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get all the matching requets
	for rows.Next() {
		i := &Invite{}

		// Pull data into request
		err := rows.Scan(&i.ID, &i.Code, &i.Username, &i.Creator, &i.Activated, &i.Created)
		if err != nil {
			return nil, err
		}

		// Add invite to collection
		invites = append(invites, i)
	}

	// Catch sql errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return invites, nil
}

// ValidateInvite check the valitity of an invite code
func (db *DB) ValidateInvite(inviteCode string) (bool, error) {

	// Used or not
	var status bool

	// Get id and password hash for given username
	row := db.QueryRow("SELECT status FROM invites WHERE code = $1", inviteCode)

	err := row.Scan(&status)
	if err == sql.ErrNoRows {
		return true, nil
	} else if err != nil {
		return true, nil
	}

	return status, nil
}

// CreateInvite add a new invite
func (db *DB) CreateInvite(creator, code string) error {

	// Empty invite id
	var id int

	stmt := `INSERT INTO invites (code, creator, status, created) VALUES($1, $2, FALSE, timezone('utc', now())) RETURNING id`

	// Add invite
	err := db.QueryRow(stmt, code, creator).Scan(&id)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	return nil
}

// FillInvite use an invite, invalidate for future use
func (db *DB) FillInvite(username, code string) error {

	// The empty used invite
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
