package main

import (
	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/sessions"
)

// App defines the global attributes
type App struct {
	HTMLDir    string
	StaticDir  string
	BookDir    string
	YoutubeDir string
	DB         *models.DB
	Storage    *session.Session
	BookBucket string
	BookAPIKey string
	Sessions   *sessions.CookieStore
}
