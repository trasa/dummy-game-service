package main

type CallEndpointRequest struct{
	Name 	string 				`json:"name"`
	Version int					`json:"version"`
	Params 	ParamsStruct	    `json:"params"`
}

type ParamsStruct struct {
	PlayerId	int		`json:"player_id"`
	Stars		int 	`json:"stars"`
}