package main

import (
	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"
	"github.com/gorilla/sessions"
)

// App structure
type App struct {
	HTMLDir   string
	StaticDir string
	Request  *models.RequestsDB
	User	  *models.UsersDB
	Sessions  *sessions.CookieStore
}
