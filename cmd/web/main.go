package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"encoding/gob"
	"github.com/Mr-Schneider/louieslibrary/pkg/models"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Cookie Store for user data
var session_store *sessions.CookieStore

func main() {

	// Command line flags for application settings
	env							:= flag.String("env", "prod", "environment")			
	addr 						:= flag.String("addr", ":4000", "HTTP network address")
	html_dir 				:= flag.String("html_dir", "./ui/html", "Path to HTML templates")
	static_dir 			:= flag.String("static_dir", "./ui/static", "Path to static assets")
	youtube_dir 		:= flag.String("youtube_dir", "./assets/youtube", "Path to youtube assets")
	dsn 						:= flag.String("dsn", "postgres://", "Postgres DSN")
	storage_server 	:= flag.String("storage_server", "http://", "s3 storage endpoint")
	book_bucket 		:= flag.String("book_bucket", "library", "bucket for book storage")
	book_api_key 		:= flag.String("book_api_key", "", "api key for google books api")
	storage_key 		:= flag.String("storage_key", "key", "s3 access key")
	storage_secret 	:= flag.String("storage_secret", "secret", "s3 access secret")
	tlsCert 				:= flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey 					:= flag.String("tls-key", "./tls/key.pem", "Path to TLS key")

	flag.Parse()

	// Database connection
	db := ConnectDB(*dsn)
	defer db.Close()

	// s3 storage connection
	storage := ConnectStorage(*storage_server, *storage_key, *storage_secret)

	// Initalize session manager
	session_store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Register user type for storing in sessions
	gob.Register(&models.User{})

	// Application instance
	app := &App{
		HTMLDir:		*html_dir,
		StaticDir: 	*static_dir,
		YoutubeDir:	*youtube_dir,
		DB:        	&models.DB{db},
		Storage:		storage,
		BookBucket: *book_bucket,
		BookAPIKey: *book_api_key,
		Sessions:  	session_store,
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

// ConnectDB
// Test connection to the db
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

// ConnectStorage
// Create a connection to the s3 server
func ConnectStorage(url, key, secret string) *session.Session {
	// Configure s3 remote
	storage_config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(url),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	// Return new s3 session for starting connections
	return session.New(storage_config)
}