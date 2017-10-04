package main

type Player struct {
	Id        	int				`json:"id,omitempty"`
	Stars      	int				`json:"stars,omitempty"`
}

type Players []Player
var players Players
var playerIndex int