package models

import (
	"database/sql"
	"time"

	"gopkg.in/guregu/null.v4"
)

// DB hold db connection
type DB struct {
	*sql.DB
}

// Request describe the request structure
type Request struct {
	ID        int
	Requester string
	Title     string
	Status    string
	BookID    string
	Created   time.Time
}

// Requests multiple requests
type Requests []*Request

// User describe the user structure
type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Role           string
	Created        time.Time
}

// Users multiple users
type Users []*User

// Invite describe the invite structure
type Invite struct {
	ID        string
	Username  null.String
	Code      string
	Creator   string
	Activated string
	Created   time.Time
}

// Invites multiple invites
type Invites []*Invite

// Book describe the book structure
type Book struct {
	ID             string
	VolumeID       string
	Title          string
	Subtitle       string
	Publisher      string
	PublishedDate  string
	PageCount      string
	MaturityRating string
	Authors        string
	Categories     string
	Description    string
	Uploader       string
	Price          string
	ISBN10         string
	ISBN13         string
	ImageLink      string
	Downloads      int
	Collected      bool
	Created        time.Time
}

// Books multiple books
type Books []*Book

// Review describe the review structure
type Review struct {
	ID       int
	BookID   string
	Username string
	Rating   string
	Review   string
	Created  time.Time
}

// Reviews multiple reviews
type Reviews []*Review

// Message describe the message structure
type Message struct {
	ID      int
	Sender  string
	Reciver string
	Read    bool
	Content string
	Created time.Time
}

// Messages multiple messages
type Messages []*Message
