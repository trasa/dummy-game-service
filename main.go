package main

import (
	"log"
	"net/http"
)

var logchan = make(chan string)

func main() {
	router := NewRouter()

	log.Println("Listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
