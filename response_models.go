package main

type ApiResponse struct {
	Success		    bool		`json:"success"`
	ResultCode	    string		`json:"result_code"`
	ResultMessage   string      `json:"result_message,omitempty"`
}

type CreatePlayerResponse struct {
	ApiResponse
	PlayerId    int `json:"player_id"`
}

type GetStarsResponse struct {
	ApiResponse
	Stars    int `json:"stars"`
}

type PlayerStatusResponse struct {
	ApiResponse
	Player      `json:"player"`
}

type TypeMap map[string]string
type PropertyMap map[string]TypeMap

type RequestResponseSchema struct{
	AdditionalProperties        	bool			    `json:"additionalProperties"`
	Required                        []string                    `json:"required"`
	Type                            string                      `json:"type"`
	Properties                      PropertyMap                 `json:"properties"`
}

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