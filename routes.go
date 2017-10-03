package main

import (
	"net/http"
)

type Param struct{
	Name 		string
	Required 	bool
	Type		string
}

type ResponseParam struct{
	Name 		string
	Type		string
}

type Route struct {
	Name        					string
	Description						string
	Method      					string				`json:method`
	Version							int
	Pattern     					string				`json:route`
	AdditionalProperties			bool
	HandlerFunc http.HandlerFunc
	RequestParams		[]Param
	ResponseParams		[]Param
}
type Routes []Route
var routes Routes
var internalRoutes Routes

func init(){

	routes = Routes{
		Route{
			Name: "Index",
			Method: "GET",
			Pattern: "/",
			HandlerFunc: Index,
		},
		Route{
			Name: "PostWebhook",
			Method: "POST",
			Pattern: "/webhook",
			HandlerFunc: PostWebhook,
		},
		Route{
			Name: "GetWebhook",
			Method: "GET",
			Pattern: "/webhook",
			HandlerFunc: GetWebhook,
		},
		Route{
			Name: "ViewLog",
			Method: "GET",
			Pattern: "/viewlog",
			HandlerFunc: LogSocket,
		},
		Route{
			Name: "Ping",
			Method: "GET",
			Pattern: "/ping",
			HandlerFunc: Ping,
		},
		Route{
			Name: "",
			Description: "Default directory",
			Method:	"GET",
			Version: 1,
			Pattern: "/",
			AdditionalProperties: true,
			HandlerFunc: Index,
			RequestParams: []Param{},
			ResponseParams:	[]Param{
			},
		},
		Route{
			Name: "gg",
			Description: "Calls specified endpoint using passed params",
			Method:	"POST",
			Version: 1,
			Pattern: "/gg",
			AdditionalProperties: true,
			HandlerFunc: CallEndpoint,
			RequestParams:	[]Param{
				Param{"name", true, "string"},
				Param{"version", false, "integer"},
				Param{"params", true, "object"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
			},
		},
		Route{
			Name: "discovery",
			Description: "Returns a list of available endpoints on the server.",
			Method:	"GET",
			Version: 1,
			Pattern: "/discovery",
			AdditionalProperties: true,
			HandlerFunc: ViewEndpoints,
			RequestParams:	[]Param{},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
			},
		},
	}

	//	these aren't actual API endpoints -- just formatted very similarly.
	//	TODO: refactor code to remove extraneous information
	internalRoutes = Routes{
		Route{
			Name:	"createplayer",
			Description:	"Creates a player on the server.",
			Version:	1,
			AdditionalProperties:	true,
			RequestParams:	[]Param{

			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"player_id", true, "integer"},
			},
		},

		Route{
			Name:	"getstars",
			Description:	"Returns the number of stars associated with specified player",
			Version:	1,
			AdditionalProperties:	true,
			RequestParams:	[]Param{
				Param{"player_id", true, "integer"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"stars", true, "integer"},
			},
		},

		Route{
			Name:	"addstars",
			Description:	"Grants a specified number of stars to the player, and displays current stars.",
			Version:	1,
			AdditionalProperties:	true,
			RequestParams:	[]Param{
				Param{"player_id", true, "integer"},
				Param{"stars", true, "integer"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"stars", true, "integer"},
			},
		},

		Route{
			Name:	"addstars",
			Description:	"Grants a specified number of stars to the player, and displays player info.",
			Version:	2,
			AdditionalProperties:	true,
			RequestParams:	[]Param{
				Param{"player_id", true, "integer"},
				Param{"stars", true, "integer"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"player", true, "object"},
			},
		},

		Route{
			Name:	"subtractstars",
			Description:	"Removes a specified number of stars from the player.",
			Version:	1,
			AdditionalProperties:	true,
			RequestParams:	[]Param{
				Param{"player_id", true, "integer"},
				Param{"stars", true, "integer"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"stars", true, "integer"},
			},
		},

		Route{
			Name:	"wipestars",
			Description:	"Removes all stars from the player. GG.",
			Version:	1,
			AdditionalProperties:	true,
			RequestParams:	[]Param{
				Param{"player_id", true, "integer"},
			},
			ResponseParams:	[]Param{
				Param{"success", true, "boolean"},
				Param{"result_code", true, "string"},
				Param{"result_message", false, "string"},
				Param{"stars", true, "integer"},
			},
		},
	}
}

