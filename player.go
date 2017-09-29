package main

type Player struct {
	Id        	int				`json:"id"`
	Stars      	int				`json:"stars"`
}

type Players []Player
var players Players
var playerIndex int