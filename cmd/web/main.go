package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// Cookie Store for user data
var sessionStore *sessions.CookieStore

func main() {

	// Command line flags for application settings
	env := flag.String("env", "prod", "environment")
	addr := flag.String("addr", ":4000", "HTTP network address")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	youtubeDir := flag.String("youtube-dir", "./assets/youtube", "Path to youtube assets")
	dsn := flag.String("dsn", "postgres://", "Postgres DSN")
	storageServer := flag.String("storage-server", "http://", "s3 storage endpoint")
	bookBucket := flag.String("book-bucket", "library", "bucket for book storage")
	bookAPIKey := flag.String("book-api-key", "", "api key for google books api")
	storageKey := flag.String("storage-key", "key", "s3 access key")
	storageSecret := flag.String("storage-secret", "secret", "s3 access secret")
	jwtKey := flag.String("jwt-key", "supersecure", "JWT secure string")
	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS key")

	flag.Parse()

	// Database connection
	db := ConnectDB(*dsn)
	defer db.Close()

	// s3 storage connection
	storage := ConnectStorage(*storageServer, *storageKey, *storageSecret)

	// Initalize session manager
	sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Register user type for storing in sessions
	gob.Register(&models.User{})

	// Application instance
	app := &App{
		HTMLDir:      *htmlDir,
		StaticDir:    *staticDir,
		YoutubeDir:   *youtubeDir,
		DB:           &models.DB{db},
		Storage:      storage,
		BookBucket:   *bookBucket,
		BookAPIKey:   *bookAPIKey,
		Sessions:     sessionStore,
		SecureString: []byte(*jwtKey),
	}

	//Start server, quit on failure
	log.Printf("Starting server on %s", *addr)
	if *env == "test" {
		err := http.ListenAndServe(*addr, app.Routes())
		log.Fatal(err)
	} else {
		err := http.ListenAndServeTLS(*addr, *tlsCert, *tlsKey, app.Routes())
		log.Fatal(err)
	}
}

// ConnectDB Test connection to the db
func ConnectDB(dsn string) *sql.DB {
	// Postgres
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Return db connection
	return db
}

// ConnectStorage Create a connection to the s3 server
func ConnectStorage(url, key, secret string) *session.Session {
	// Configure s3 remote
	storageConfig := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(url),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Return new s3 session for starting connections
	return session.New(storageConfig)
}
