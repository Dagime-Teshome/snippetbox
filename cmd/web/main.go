package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Dagime-Teshome/snippetbox/pkg/models"
	"github.com/Dagime-Teshome/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippet  interface {
		Insert(string, string, string) (int, error)
		GetById(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	user interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	templateCache map[string]*template.Template
	Session       *sessions.Session
}

func main() {

	infoLog := log.New(os.Stdout, "Info:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// data that a user can pass when starting the web server (data from the command line)
	dsn := flag.String("dsn", "root:root@/snippetbox?parseTime=true", "connection string for database")
	addr := flag.String("addr", ":3001", "port for server")
	secret := flag.String("secret", "xK7vQp3Lz9R4mN8tY2cHf6Bq1Ws5Dj0U", "secret to encrypt session")
	flag.Parse()

	// Open Database connection
	db, err := opeDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// initialize template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// setup and create session
	session := sessions.New([]byte(*secret))
	session.Lifetime = time.Minute * 30
	session.SameSite = http.SameSiteStrictMode
	session.Secure = true

	// injected dependency of the app that handlers and middleware use
	app := application{

		errorLog:      errorLog,
		infoLog:       infoLog,
		user:          &mysql.UserModel{Db: db},
		snippet:       &mysql.SnippetModel{Db: db},
		templateCache: templateCache,
		Session:       session,
	}
	// info text outputted at server start
	infoLog.Printf("Starting server on %s", *addr)

	// select cypher type for tls (makes message sending secure)
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// setting for the server
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
		errorLog.Fatal("server:%w", err)
	}
}

func opeDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
