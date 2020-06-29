package main

import (
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
	"github.com/gorilla/sessions"
	"github.com/aws/aws-sdk-go/aws/session"
)

// App structure
type App struct {
	HTMLDir   string
	StaticDir string
	BookDir   string
	YoutubeDir string
	DB        *models.DB
	Storage		*session.Session
	Sessions  *sessions.CookieStore
}
