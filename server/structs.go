package main

type Player struct {
	Token string
	Name  string
	RLFP  PlayerFP
}

type PlayerPosition struct {
	Name string
	X    float32
	Y    float32
	Z    float32
}
