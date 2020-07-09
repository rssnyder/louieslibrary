package models

import (
	"time"
	"database/sql"
	"gopkg.in/guregu/null.v4"
)

type DB struct {
	*sql.DB
}

// Request
// Describe the request structure
type Request struct {
	ID        int
	Requester string
	Title     string
	Status    string
	BookID    string
	Created   time.Time
}

// Requests
type Requests []*Request

// User
// Describe the user structure
type User struct {
	ID        			int
	Username  			string
	Email     			string
	HashedPassword	[]byte
	Role      			string
	Created   			time.Time
}

// Users
type Users []*User

// Invite
// Describe the invite structure
type Invite struct {
	ID					string
	Username 		null.String
	Code    		string
	Creator			string
	Activated		string
	Created   	time.Time
}

// Invites
type Invites []*Invite

// Book
// Describe the book structure
type Book struct {
	ID							string
	VolumeID				string
	Title       		string
	Subtitle				string
	Publisher				string
	PublishedDate		string
	PageCount				string
	MaturityRating	string
	Authors      		string
	Categories      string
	Description 		string
	Uploader 				string
	Price						string
	ISBN10					string
	ISBN13					string
	ImageLink				string
	Downloads				int
	Collected				bool
	Created    			time.Time
}

// Books
type Books []*Book

// Review
// Describe the review structure
type Review struct {
	ID        int
	BookID    string
	Username	string
	Rating    string
	Review    string
	Created   time.Time
}

// Reviews
type Reviews []*Review

// Message
// Describe the message structure
type Message struct {
	ID      int
	Sender  string
	Reciver	string
	Read    bool
	Content	string
	Created	time.Time
}

// Messages
type Messages []*Message