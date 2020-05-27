package forms

import (
	"strings"
	"unicode/utf8"
)

// NewRequest models the request structure
type NewRequest struct {
	Requester string
	Title     string
	Failures  map[string]string
}

// Valid makes sure the the fields are correctly formatted
func (f *NewRequest) Valid() bool {
	f.Failures = make(map[string]string)

	// Check for non-empty Requester
	if strings.TrimSpace(f.Requester) == "" {
		f.Failures["Requester"] = "Requester is required"
	} else if utf8.RuneCountInString(f.Requester) > 100 {
		f.Failures["Requester"] = "Requester cannot be longer than 100 characters"
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
	}

	return len(f.Failures) == 0
}
