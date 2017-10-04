package main

type ApiResponse struct {
	Success		    bool		`json:"success"`
	ResultCode	    string		`json:"result_code"`
	ResultMessage   string      `json:"result_message,omitempty"`
}

//type CreatePlayerResponse struct {
//	ApiResponse
//	PlayerId    int `json:"player_id"`
//}
//
//type GetStarsResponse struct {
//	ApiResponse
//	Stars    int `json:"stars"`
//}
//
//type PlayerStatusResponse struct {
//	ApiResponse
//	Player      `json:"player"`
//}

type GenericResponse struct{
	ApiResponse
	Properties	`json:"properties,omitempty"`
}

type Properties struct{
	*Player			`json:"player,omitempty"`
	PlayerId	*int	`json:"player_id,omitempty"`
	Stars		*int	`json:"stars,omitempty"`
}

type TypeMap map[string]string

type RequestResponseSchema struct{
	SchemaRoot
	SchemaBody
}

type SchemaRoot struct{
	SchemaUrl string `json:"$schema,omitempty"`
	Definitions map[string]string `json:"definitions,omitempty"`
}

// a schema that describes a property
type SchemaBody struct{
	Type                            string                      `json:"type"`
	AdditionalProperties        	bool			    `json:"additionalProperties,omitempty"`
	Required                        []string                    `json:"required,omitempty"`
	Properties                      PropertyMap                 `json:"properties,omitempty"`
}

type PropertyMap map[string]SchemaBody

type GGMethodSchema struct{
	Name        string  `json:"method_name"`
	Description string  `json:"description"`
	Version     int     `json:"version"`

	RequestSchema       RequestResponseSchema   `json:"request_schema"`
	ResponseSchema      RequestResponseSchema   `json:"response_schema"`
}

type DiscoveryResponse struct{
	ApiResponse
	Methods []GGMethodSchema  `json:"methods"`
}