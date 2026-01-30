package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Write([]byte("404 not found"))
		return
	}
	w.Write([]byte("Home page"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
		return
	}
	w.Write([]byte("create snippet"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if r.URL.Query().Get("id") == "" {
			id = 0 // Default value
			fmt.Printf("Using default value: %d\n", id)
		}
		http.Error(w, "id must be a valid number", 505)
		return
	}
	if id <= 0 {
		http.Error(w, "id needs to be greater than zero", 505)
		return
	}

	w.Write([]byte(fmt.Sprintf("show snippet for %v", id)))
}

func main() {
	fmt.Println("hello")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	fmt.Println("Starting server on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("server:%w", err)
	}
}
