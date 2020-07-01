package models

import (
	"time"
	"database/sql"
	"gopkg.in/guregu/null.v4"
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
	BookID    string
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

// Invite holds invite codes
type Invite struct {
	ID					string
	Username 		null.String
	Code    		string
	Creator			string
	Activated		string
	Created   	time.Time
}

// Invites holds a users invites
type Invites []*Invite

// Book holds data on a book
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
	Created    			time.Time
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