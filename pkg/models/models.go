package models

import (
	"time"
)

// Request describes the Request structure
type Request struct {
	ID        int
	Requester string
	Title     string
	Completed bool
	Created   time.Time
}

// Requests holds multiple Requests
type Requests []*Request
