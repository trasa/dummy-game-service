package main

import (
	"fmt"
	"log"
	"html"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you are at %q", html.EscapeString(r.URL.Path))
}

func PostWebhook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("Posted Form Data: %v", r.Form)
	fmt.Fprintf(w, "POST Received")
}

func GetWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET received")
}
