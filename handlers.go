package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"io/ioutil"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/#", 302)
}

func PostWebhook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("Posted Form Data:\n %v\n%v", r.Header, r.Form)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading form body: %s", err)
	}
	broadcaster <- fmt.Sprintf("Webhook Form Received:\nHeaders: %v\nForm: %v\nBody: %s\n", r.Header, r.Form, body)

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
		fmt.Fprintln(w, "POST Received")
	}
}

func GetWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET received")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PONG")
	broadcaster <- "PONG"
}

func LogSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	client := newClient(c)
	clients.Add(client)
	go client.writePump()
	client.readPump()
}
