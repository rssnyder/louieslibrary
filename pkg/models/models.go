package models

import (
	"time"
)

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

// User hold data on a logged in user
type User struct {
	ID        int
	Username  string
	Email     string
	HashedPassword []byte
	Role      string
	Created   time.Time
}