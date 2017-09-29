package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

//	calls the endpoint with associated params, specified inside the body
func CallEndpoint(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	functionName, reqVersion, params, err := ParseRequest(w, r)
	if err != nil {
		SendErrorResponse(w, "Deserialization error while parsing input.")
		log.Printf("Deserialization error while parsing input: %s", err)
		return
	}
	switch functionName {
	case "createplayer":
		CreatePlayer(w)
	case "getstars":
		GetStars(w, params.PlayerId)
	case "addstars":
		if reqVersion == 2{
			//	only difference is that "player":{"id":..., "stars":...} gets displayed instead of just "stars":...
			AddStarsV2(w, params)
		}	else	{
			AlterStars(w, params, false)
		}
	case "subtractstars":
		params.Stars = params.Stars*(-1)
		AlterStars(w, params, false)
	case "wipestars":
		AlterStars(w, params, true)
	default:
		response := ApiResponse{
			Success: false,
			ResultCode: "ERROR",
			ResultMessage: fmt.Sprint("Unrecognized method call: ", functionName),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error while calling endpoint: %s", err)
		}

	}
}

//	returns the `name` and `params` fields out of passed request
//	if error encountered during deserialization, sends back error.
func ParseRequest(w http.ResponseWriter, r *http.Request) (string, int, ParamsStruct, error){
	decoder := json.NewDecoder(r.Body)
	var t CallEndpointRequest
	err := decoder.Decode(&t)
	return t.Name, t.Version, t.Params, err
}

//	returns a list of callable endpoints on the server, along with associated metadata
func ViewEndpoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	response := DiscoveryResponse{
		ApiResponse:	ApiResponse{Success: true, ResultCode: "OK"},
		Methods:		[]GGMethodSchema{},
	}

	for _, route := range internalRoutes{
		ggMethodSchema := GGMethodSchema{
			Name:			route.Name,
			Description:	route.Description,
			Version:		route.Version,
			RequestSchema:	CreateRequestResponseSchema(route, route.RequestParams),
			ResponseSchema:	CreateRequestResponseSchema(route, route.ResponseParams),
		}

		response.Methods = append(response.Methods, ggMethodSchema)

	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while requesting endpoints: %s", err)
	}
}

//	creates a request/response schema populated with passed route and param information
func CreateRequestResponseSchema(route Route, params []Param) RequestResponseSchema{
	requestResponseSchema := RequestResponseSchema{
		AdditionalProperties: 	route.AdditionalProperties,
		Required:				[]string{},
		Type:					"object",
		Properties:				PropertyMap{},
	}

	for _, parameter := range params{
		typeMap := TypeMap{"type": parameter.Type}
		requestResponseSchema.Properties[parameter.Name] = typeMap
		if parameter.Required == true{
			requestResponseSchema.Required = append(requestResponseSchema.Required, parameter.Name)
		}
	}

	return requestResponseSchema
}

//	creates a new player and increments the global index
func CreatePlayer(w http.ResponseWriter) {

	playerIndex+=1
	newPlayer := Player{playerIndex, 0}
	players = append(players, newPlayer)

	apiResponse := ApiResponse{Success: true, ResultCode: "OK"}
	response := CreatePlayerResponse{
		apiResponse,
		newPlayer.Id,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while creating player: %s", err)
	}
}

//	returns response with number of stars belonging to specified player
func GetStars(w http.ResponseWriter, playerId int){
	//	validate input
	if playerId <= 0 {
		SendErrorResponse(w, fmt.Sprint("Illegal value for 'player_id': ", playerId))
		return
	}

	playerPtr := FindPlayerById(playerId)
	if playerPtr != nil {
		SendSuccessGetStarsResponse(w, playerPtr)
	}	else	{
		SendErrorResponse(w, fmt.Sprint("Unable to find player with id: ", playerId))
	}
}

//	alters a player's star count by adding the stars in the ParamsStruct to them (could be a negative number)
//	 or wiping them to zero, based on the bool.
func AlterStars(w http.ResponseWriter, params ParamsStruct, doWipe bool){
	//	validate input
	if params.PlayerId <= 0 {
		SendErrorResponse(w, fmt.Sprint("Illegal value for 'player_id': ", params.PlayerId))
		return
	}
	playerPtr := FindPlayerById(params.PlayerId)

	if playerPtr != nil {

		(*playerPtr).Stars += params.Stars
		//	if balance dips below or wipe requested, set to zero.
		if (doWipe == true) || ((*playerPtr).Stars < 0){
			(*playerPtr).Stars = 0
		}

		SendSuccessGetStarsResponse(w, playerPtr)
	}	else	{
		SendErrorResponse(w, fmt.Sprint("Unable to find player with id: ", params.PlayerId))
	}
}

//
func AddStarsV2(w http.ResponseWriter, params ParamsStruct){
	//	validate input
	if params.PlayerId <= 0 {
		SendErrorResponse(w, fmt.Sprint("Illegal value for 'player_id': ", params.PlayerId))
		return
	}

	playerPtr := FindPlayerById(params.PlayerId)

	if playerPtr != nil {
		(*playerPtr).Stars += params.Stars
		SendSuccessPlayerStatusResponse(w, playerPtr)
	}	else	{
		SendErrorResponse(w, fmt.Sprint("Unable to find player with id: ", params.PlayerId))
	}
}

//	sends a generic `OK` response along with player information
func SendSuccessPlayerStatusResponse(w http.ResponseWriter, playerPtr *Player){
	apiResponse := ApiResponse{Success: true, ResultCode: "OK"}
	response := PlayerStatusResponse{
		apiResponse,
		(*playerPtr),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while serializing response: %s", err)
	}
}

//	sends a generic `OK` response along with player information
func SendSuccessGetStarsResponse(w http.ResponseWriter, playerPtr *Player){
	apiResponse := ApiResponse{Success: true, ResultCode: "OK"}
	response := GetStarsResponse{
		apiResponse,
		(*playerPtr).Stars,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while serializing response: %s", err)
	}
}

//	sends generic ERROR response
func SendErrorResponse(w http.ResponseWriter, message string){
	//	TODO: note source of error prior to calling function
	apiResponse := ApiResponse{
		Success: false,
		ResultCode: "ERROR",
		ResultMessage: message,
	}
	if err := json.NewEncoder(w).Encode(apiResponse); err != nil {
		log.Printf("Error while serializing error response: %s", err)
	}
}

//	return max of two numbers
func Max(a, b int) int {
	if a > b{
		return a
	}
	return b
}

//	look for player with specified id inside global `players` list
func FindPlayerById(playerId int) *Player{
	var ptr *Player
	for i, player := range players{
		if player.Id == playerId{
			ptr = &players[i]
			break
		}
	}
	return ptr
}

//	attempt to extract `playerId` and `stars` from incoming request
func ExtractParamsFromBody(r *http.Request) (int, int){
	decoder := json.NewDecoder(r.Body)
	var t CallEndpointRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Printf("Error while extracting input params: %s", err)
	}
	return t.Params.PlayerId, t.Params.Stars
}
