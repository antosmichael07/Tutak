package main

type Player struct {
	Token string
	Name  string
	RLFP  PlayerFP
}

type PlayerPositionToSend struct {
	Name string
	X    float32
	Y    float32
	Z    float32
}

type PlayerRotationToSend struct {
	Name string
	X    float32
	Y    float32
}
