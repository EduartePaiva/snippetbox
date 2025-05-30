package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static"))))

	log.Println("Starting server at port: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
