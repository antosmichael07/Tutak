package main

import rl "github.com/gen2brain/raylib-go/raylib"

func new_block_model(image *rl.Image) rl.Model {
	cube_mesh := rl.GenMeshCube(1., 1., 1.)
	model := rl.LoadModelFromMesh(cube_mesh)
	rl.ImageFlipVertical(image)
	texture := rl.LoadTextureFromImage(image)
	model.GetMaterials()[0].Maps.Texture = texture

	rl.ImageFlipVertical(image)

	return model
}

func draw_block(model rl.Model, position rl.Vector3) {
	rl.DrawModelEx(model, position, rl.NewVector3(0., 0., 0.), 0., rl.NewVector3(1., 1., 1.), rl.White)
}

func new_block(name string, image *rl.Image) Block {
	block_model := new_block_model(image)

	return Block{name, block_model, false, nil}
}

func (block *Block) on_interact(fn func()) {
	block.OnInteract = fn
	block.Interactable = true
}
