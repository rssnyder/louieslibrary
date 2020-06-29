package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"encoding/gob"

	"github.com/Mr-Schneider/request.thecornelius.duckdns.org/pkg/models"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var session_store *sessions.CookieStore

func main() {
	// Flags
	env							:= flag.String("env", "prod", "environment")			
	addr 						:= flag.String("addr", ":4000", "HTTP network address")
	html_dir 				:= flag.String("html_dir", "./ui/html", "Path to HTML templates")
	static_dir 			:= flag.String("static_dir", "./ui/static", "Path to static assets")
	book_dir				:= flag.String("book_dir", "./assets/books", "Path to book assets")
	youtube_dir 		:= flag.String("youtube_dir", "./assets/youtube", "Path to youtube assets")
	dsn 						:= flag.String("dsn", "postgres://", "Postgres DSN")
	storage_server 	:= flag.String("storage_server", "http://", "s3 storage endpoint")
	storage_key 		:= flag.String("storage_key", "key", "s3 access key")
	storage_secret 	:= flag.String("storage_secret", "secret", "s3 access secret")
	tlsCert 				:= flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey 					:= flag.String("tls-key", "./tls/key.pem", "Path to TLS key")


	flag.Parse()

	// Database connection
	db := connect_db(*dsn)
	defer db.Close()

	// s3 storage connection
	storage := connect_storage(*storage_server, *storage_key, *storage_secret)

	// Initalize session manager
	session_store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Register user type for storing in sessions
	gob.Register(&models.User{})

	// Application instance
	app := &App{
		HTMLDir:		*html_dir,
		StaticDir: 	*static_dir,
		BookDir:   	*book_dir,
		YoutubeDir:	*youtube_dir,
		DB:        	&models.DB{db},
		Storage:		storage,
		Sessions:  	session_store,
	}

	//app.UploadObject("dump", "test")
	//app.DownloadObject("dump", "test", "thisismyminiodownload")
	app.DownloadBytes("library", "p - p.mobi")

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

// connect DB connection setup
func connect_db(dsn string) *sql.DB {
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

	return db
}

func connect_storage(url, key, secret string) *session.Session {
	// Configure s3 remote
	storage_config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(url),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	return session.New(storage_config)
}