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

// NewUser models the user signup
type NewUser struct {
	Username string
	Email string
	Password     string
	Failures  map[string]string
}

// Valid makes sure the the fields are correctly formatted
func (f *NewUser) Valid() bool {
	f.Failures = make(map[string]string)

	// Check for non-empty username
	if strings.TrimSpace(f.Username) == "" {
		f.Failures["Username"] = "Username is required"
	} else if utf8.RuneCountInString(f.Username) > 60 {
		f.Failures["Username"] = "Username cannot be longer than 60 characters"
	}

	// Check for non-empty email
	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	}

	// Check for non-empty password
	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
	} else if utf8.RuneCountInString(f.Password) < 8 {
		f.Failures["Password"] = "Password cannot be less than 8 characters"
	}

	return len(f.Failures) == 0
}