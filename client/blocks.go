package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	GRASS_BLOCK = iota
)

var block_map = []Block{
	new_block("Grass Block", rl.LoadImage("assets/grass_block.png")),
}
