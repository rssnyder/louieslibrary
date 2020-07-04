package forms

import (
	"strings"
	"unicode/utf8"
	"log"
)

// NewRequest models the request structure
type NewRequest struct {
	Requester	string
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
	Username 		string
	Email    		string
	Password 		string
	InviteCode	string
	Failures 		map[string]string
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

	// Check for non-empty invite code
	if strings.TrimSpace(f.InviteCode) == "" {
		f.Failures["InviteCode"] = "InviteCode is required"
		log.Printf("User submitted with InviteCode missing")
	} else if utf8.RuneCountInString(f.InviteCode) > 36 {
		f.Failures["InviteCode"] = "InviteCode cannot be more than 36 characters"
		log.Printf("User submitted with InviteCode over limit")
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
	ID 							string
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
	Failures    		map[string]string
}

// Valid makes sure the the fields are correctly formatted
func (f *NewBook) Valid() bool {
	f.Failures = make(map[string]string)

	// Check for non-empty VolumeID
	if strings.TrimSpace(f.VolumeID) == "" {
		f.Failures["VolumeID"] = "VolumeID is required"
		log.Printf("Book submitted with no VolumeID")
	}
	
	// Check for non-empty Title
	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required"
		log.Printf("Book submitted with no Title")
	}

	// Check for non-empty Subtitle
	if strings.TrimSpace(f.Subtitle) == "" {
		f.Failures["Subtitle"] = "Subtitle is required"
		log.Printf("Book submitted with no Subtitle")
	}

	// Check for non-empty Publisher
	if strings.TrimSpace(f.Publisher) == "" {
		f.Failures["Publisher"] = "Publisher is required"
		log.Printf("Book submitted with no Publisher")
	}

	// Check for non-empty PublishedDate
	if strings.TrimSpace(f.PublishedDate) == "" {
		f.Failures["PublishedDate"] = "PublishedDate is required"
		log.Printf("Book submitted with no PublishedDate")
	} else if utf8.RuneCountInString(f.PublishedDate) > 50 {
		f.Failures["PublishedDate"] = "PublishedDate cannot be longer than 50 characters"
		log.Printf("Book submitted with PublishedDate over limit")
	}

	// Check for non-empty PageCount
	if strings.TrimSpace(f.PageCount) == "" {
		f.Failures["PageCount"] = "PageCount is required"
		log.Printf("Book submitted with no PageCount")
	} else if utf8.RuneCountInString(f.PageCount) > 10 {
		f.Failures["PageCount"] = "PageCount cannot be longer than 10 characters"
		log.Printf("Book submitted with PageCount over limit")
	}

	// Check for non-empty MaturityRating
	if strings.TrimSpace(f.MaturityRating) == "" {
		f.Failures["MaturityRating"] = "MaturityRating is required"
		log.Printf("Book submitted with no MaturityRating")
	}

	// Check for non-empty Authors
	if strings.TrimSpace(f.Authors) == "" {
		f.Failures["Authors"] = "Authors is required"
		log.Printf("Book submitted with no Authors")
	}

	// Check for non-empty Categories
	if strings.TrimSpace(f.Categories) == "" {
		f.Failures["Categories"] = "Categories is required"
		log.Printf("Book submitted with no Categories")
	}

	// Check for non-empty Description
	if strings.TrimSpace(f.Description) == "" {
		f.Failures["Description"] = "Description is required"
		log.Printf("Book submitted with no Description")
	}

	// Check for non-empty Uploader
	if strings.TrimSpace(f.Uploader) == "" {
		f.Failures["Uploader"] = "Uploader is required"
		log.Printf("Book submitted with no Uploader")
	} else if utf8.RuneCountInString(f.Uploader) > 50 {
		f.Failures["Uploader"] = "Uploader cannot be longer than 50 characters"
		log.Printf("Book submitted with Uploader over limit")
	}

	// Check for non-empty Price
	if strings.TrimSpace(f.Price) == "" {
		f.Failures["Price"] = "Price is required"
		log.Printf("Book submitted with no Price")
	} else if utf8.RuneCountInString(f.Price) > 10 {
		f.Failures["Price"] = "Price cannot be longer than 10 characters"
		log.Printf("Book submitted with Price over limit")
	}

	// Check for non-empty ISBN10
	if strings.TrimSpace(f.ISBN10) == "" {
		f.Failures["ISBN10"] = "ISBN10 is required"
		log.Printf("Book submitted with no ISBN10")
	} else if utf8.RuneCountInString(f.ISBN10) > 10 {
		f.Failures["ISBN10"] = "ISBN10 cannot be longer than 10 characters"
		log.Printf("Book submitted with ISBN10 over limit")
	}

	// Check for non-empty ISBN13
	if strings.TrimSpace(f.ISBN13) == "" {
		f.Failures["ISBN13"] = "ISBN13 is required"
		log.Printf("Book submitted with no ISBN13")
	} else if utf8.RuneCountInString(f.ISBN13) > 13 {
		f.Failures["ISBN13"] = "ISBN13 cannot be longer than 13 characters"
		log.Printf("Book submitted with ISBN13 over limit")
	}

	// Check for non-empty ImageLink
	if strings.TrimSpace(f.ImageLink) == "" {
		f.Failures["ImageLink"] = "ImageLink is required"
		log.Printf("Book submitted with no ImageLink")
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