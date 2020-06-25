package models

import (
	"time"
	"database/sql"
)

type DB struct {
	*sql.DB
}

// Request describes the Request structure
type Request struct {
	ID        int
	Requester string
	Title     string
	Status    string
	Created   time.Time
}

// Requests holds multiple Requests
type Requests []*Request

// User holds data on a logged in user
type User struct {
	ID        int
	Username  string
	Email     string
	HashedPassword []byte
	Role      string
	Created   time.Time
}

// Book holds data on a book
type Book struct {
	ID         int
	ISBN       string
	Title      string
	Author     string
	Genre      string
	Description string
	Uploader    string
	Created    time.Time
}

// Books holds multiple books
type Books []*Book

// Review holds a review
type Review struct {
	ID         int
	BookID       string
	Username      string
	Rating     string
	Review      string
	Created    time.Time
}

// Books holds multiple books
type Reviews []*Review