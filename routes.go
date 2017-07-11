package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"PostWebhook",
		"POST",
		"/webhook",
		PostWebhook,
	},
	Route{
		"GetWebhook",
		"GET",
		"/webhook",
		GetWebhook,
	},
	Route{
		"ViewLog",
		"GET",
		"/viewlog",
		LogSocket,
	},
}
