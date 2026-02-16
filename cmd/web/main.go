package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	snippets "github.com/Dagime-Teshome/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippet       *snippets.SnippetModel
	templateCache map[string]*template.Template
	Session       *sessions.Session
}

func main() {

	infoLog := log.New(os.Stdout, "Info:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	dsn := flag.String("dsn", "root:root@/snippetbox?parseTime=true", "connection string for database")
	addr := flag.String("addr", ":3001", "port for server")
	secret := flag.String("secret", "xK7vQp3Lz9R4mN8tY2cHf6Bq1Ws5Dj0U", "secret to encrypt session")
	flag.Parse()
	db, err := opeDb(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = time.Minute * 30
	app := application{

		errorLog:      errorLog,
		infoLog:       infoLog,
		snippet:       &snippets.SnippetModel{Db: db},
		templateCache: templateCache,
		Session:       session,
	}

	infoLog.Printf("Starting server on %s", *addr)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
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
