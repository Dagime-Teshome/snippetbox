package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "Info:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)
	addr := flag.String("addr", ":3001", "port for server")
	flag.Parse()
	fmt.Println("hello")

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
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
