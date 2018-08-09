package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//Tile is used to make the map of reusable pieces
type Tile struct {
	//The attributes of the tile
	mBox sdl.Rect

	//The tile type
	mType int
}

//NewTile initializes position and type
func NewTile(x, y int32, tileType int) *Tile {
	return &Tile{
		//Get the offsets and set the collision box
		mBox: sdl.Rect{X: x, Y: y, W: tileWidth, H: tileHeight},
		//Get the tile type
		mType: tileType,
	}
}

//Render show the tile
func (t *Tile) Render(camera *sdl.Rect) error {
	//If the tile is on the screen
	if checkCollision(*camera, t.mBox) {
		//Show the tile
		err := gTileTexture.Render(t.mBox.X-camera.X, t.mBox.Y-camera.Y,
			&gTileClips[t.mType], 0, nil, sdl.FLIP_NONE)
		if err != nil {
			return fmt.Errorf("could not render tile's texture: %v", err)
		}
	}

	return nil
}

//MType exports the file type
func (t *Tile) MType() int {
	return t.mType
}

//MBox exports the collision box
func (t *Tile) MBox() sdl.Rect {
	return t.mBox
}
