package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	snippets "github.com/Dagime-Teshome/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippet       *snippets.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	infoLog := log.New(os.Stdout, "Info:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	dsn := flag.String("dsn", "root:root@/snippetbox?parseTime=true", "connection string for database")
	addr := flag.String("addr", ":3001", "port for server")
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
	app := application{

		errorLog:      errorLog,
		infoLog:       infoLog,
		snippet:       &snippets.SnippetModel{Db: db},
		templateCache: templateCache,
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
