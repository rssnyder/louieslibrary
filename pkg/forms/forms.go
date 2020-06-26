package forms

import (
	"strings"
	"unicode/utf8"
	"log"
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
		log.Printf("Request submitted with requester missing")
	} else if utf8.RuneCountInString(f.Requester) > 100 {
		f.Failures["Requester"] = "Requester cannot be longer than 100 characters"
		log.Printf("Request submitted with requester over limit")
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
		log.Printf("Request submitted with title missing")
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
		log.Printf("User submitted with username missing")
	} else if utf8.RuneCountInString(f.Username) > 60 {
		f.Failures["Username"] = "Username cannot be longer than 60 characters"
		log.Printf("User submitted with username over limit")
	}

	// Check for non-empty email
	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
		log.Printf("User submitted with email missing")
	}

	// Check for non-empty password
	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
		log.Printf("User submitted with password missing")
	} else if utf8.RuneCountInString(f.Password) < 8 {
		f.Failures["Password"] = "Password cannot be less than 8 characters"
		log.Printf("User submitted with password under limit")
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
		log.Printf("Book submitted with no isbn")
	} else if utf8.RuneCountInString(f.ISBN) > 20 {
		f.Failures["ISBN"] = "ISBN cannot be longer than 20 characters"
		log.Printf("Book submitted with ISBN over limit")
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
		log.Printf("Book submitted with no title")
	}

	// Check for non-empty Author
	if strings.TrimSpace(f.Author) == "" {
		f.Failures["Author"] = "Author is required"
		log.Printf("Book submitted with no author")
	} else if utf8.RuneCountInString(f.Author) > 50 {
		f.Failures["Author"] = "Author cannot be longer than 50 characters"
		log.Printf("Book submitted with author over limit")
	}

	// Check for non-empty Description
	if strings.TrimSpace(f.Description) == "" {
		f.Failures["Description"] = "Description is required"
		log.Printf("Book submitted with no description")
	}

	// Check for non-empty Genre
	if strings.TrimSpace(f.Genre) == "" {
		f.Failures["Genre"] = "Genre is required"
		log.Printf("Book submitted with no genre")
	} else if utf8.RuneCountInString(f.Genre) > 50 {
		f.Failures["Genre"] = "Genre cannot be longer than 50 characters"
		log.Printf("Book submitted with genre over limit")
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
		log.Printf("Review submitted missing book id")
	}

	// Check for non-empty title
	if strings.TrimSpace(f.Username) == "" {
		f.Failures["Username"] = "Username is required"
		log.Printf("Review submitted missing username")
	}

	// Check for non-empty Rating
	if strings.TrimSpace(f.Rating) == "" {
		f.Failures["Rating"] = "Rating is required"
		log.Printf("Review submitted missing rating")
	} else if utf8.RuneCountInString(f.Rating) > 1 {
		f.Failures["Rating"] = "Rating cannot be longer than 1 character"
		log.Printf("Review submitted with invalid rating")
	}

	// Check for non-empty review
	if strings.TrimSpace(f.Review) == "" {
		f.Failures["Review"] = "Review is required"
		log.Printf("Review submitted missing review")
	}
	return len(f.Failures) == 0
}