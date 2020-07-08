package main

import (
	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/gorilla/sessions"
	"github.com/aws/aws-sdk-go/aws/session"
)

// App
// This defines the structure of the application
// and the things it requires to operate
type App struct {
	HTMLDir   	string
	StaticDir 	string
	BookDir   	string
	YoutubeDir 	string
	DB        	*models.DB
	Storage			*session.Session
	BookBucket	string
	BookAPIKey 	string
	Sessions  	*sessions.CookieStore
}
