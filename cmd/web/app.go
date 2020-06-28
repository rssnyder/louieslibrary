package main

import (
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
	"github.com/gorilla/sessions"
)

// App structure
type App struct {
	HTMLDir   string
	StaticDir string
	BookDir   string
	YoutubeDir string
	DB        *models.DB
	Sessions  *sessions.CookieStore
}
