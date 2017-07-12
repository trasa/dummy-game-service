package main

import (
	"log"
	"net/http"
)

var globalQuit = make(chan struct{})
var clients = newClients()
var broadcaster = make(chan string, 10)

func main() {
	// dispatcher for all incoming messages
	// from broadcaster to all clients
	go func() {
		for {
			msg := <-broadcaster
			clients.Iter(func(c *Client) { c.source <- msg })
		}
	}()

	router := NewRouter()

	log.Println("Listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
	close(globalQuit)
}
