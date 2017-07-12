package main

import (
	"log"
	"net/http"
)

// this tells all clients to terminate
var globalQuit = make(chan struct{})

// our set of connected clients
var clients = newClients()

// the channel that gets dispatched to all clients
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
