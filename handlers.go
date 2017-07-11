package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you are at %q", html.EscapeString(r.URL.Path))
}

func PostWebhook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("Posted Form Data: %v", r.Form)

	// is this a fail?
	failValue := r.PostFormValue("fail")
	if failValue != "" {
		// force fail to happen
		log.Printf("forcing a fail: %s", failValue)
		responseCode, err := strconv.Atoi(failValue)
		if err != nil {
			// couldn't decode response code so use a default
			responseCode = 500
		}
		http.Error(w, "forcing an error", responseCode)
	} else {
		// otherwise success
		fmt.Fprintf(w, "POST Received")
	}
}

func GetWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET received")
}

func LogSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request to connect logview websocket")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	client := newClient(c)
	go client.writePump()
	client.readPump()
}
