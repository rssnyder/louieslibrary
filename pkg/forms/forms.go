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
	Email    string
	Password string
	Failures map[string]string
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

// Book holds data on a book
type NewBook struct {
	ID          int
	ISBN        string
	Title       string
	Author      string
	Genre       string
	Description string
	Upload		bool
	Uploader string
	Failures    map[string]string
}

// Valid makes sure the the fields are correctly formatted
func (f *NewBook) Valid() bool {
	f.Failures = make(map[string]string)

	// Check for non-empty isbn
	if strings.TrimSpace(f.ISBN) == "" {
		f.Failures["ISBN"] = "ISBN is required"
	} else if utf8.RuneCountInString(f.ISBN) > 20 {
		f.Failures["ISBN"] = "ISBN cannot be longer than 20 characters"
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
	}

	// Check for non-empty Author
	if strings.TrimSpace(f.Author) == "" {
		f.Failures["Author"] = "Author is required"
	} else if utf8.RuneCountInString(f.Author) > 50 {
		f.Failures["Author"] = "Author cannot be longer than 50 characters"
	}

	// Check for non-empty Description
	if strings.TrimSpace(f.Description) == "" {
		f.Failures["Description"] = "Description is required"
	}

	// Check for non-empty Genre
	if strings.TrimSpace(f.Genre) == "" {
		f.Failures["Genre"] = "Genre is required"
	} else if utf8.RuneCountInString(f.Genre) > 50 {
		f.Failures["Genre"] = "Genre cannot be longer than 50 characters"
	}

	return len(f.Failures) == 0
}

// Review holds data on a review
type NewReview struct {
	ID        int
	BookID    string
	Username  string
	Rating    string
	Review    string
	Failures  map[string]string
}

// Valid makes sure the the fields are correctly formatted
func (f *NewReview) Valid() bool {
	f.Failures = make(map[string]string)

	// Check for empty book id
	if strings.TrimSpace(f.BookID) == "" {
		f.Failures["BookID"] = "BookID is required"
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Username) == "" {
		f.Failures["Username"] = "Username is required"
	}

	// Check for non-empty Rating
	if strings.TrimSpace(f.Rating) == "" {
		f.Failures["Rating"] = "Rating is required"
	} else if utf8.RuneCountInString(f.Rating) > 1 {
		f.Failures["Rating"] = "Rating cannot be longer than 1 character"
	}

	// Check for non-empty review
	if strings.TrimSpace(f.Review) == "" {
		f.Failures["Review"] = "Review is required"
	}
	return len(f.Failures) == 0
}