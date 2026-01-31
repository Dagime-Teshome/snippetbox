package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello")
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	fmt.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("server:%w", err)
	}
}
