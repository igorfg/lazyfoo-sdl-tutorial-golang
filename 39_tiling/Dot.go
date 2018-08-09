package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//DotWidth is the dot's width
	DotWidth = 20
	//DotHeight is the dot's height
	DotHeight = 20

	//DotVel is the maximum axis velocity of the dot
	DotVel = 10
)

//Dot is the dot that will move around on the screen
type Dot struct {
	//Collision box of the dot
	mBox sdl.Rect

	//The velocity of the dot
	mVelX, mVelY int32
}

//NewDot initializes a dot
func NewDot() *Dot {
	//Initialize collision box and velocity
	return &Dot{
		mBox:  sdl.Rect{X: 0, Y: 0, W: DotHeight, H: DotWidth},
		mVelX: 0,
		mVelY: 0,
	}
}

//HandleEvent takes keypresses and adjusts the dot's velocity
func (d *Dot) HandleEvent(e sdl.Event) {
	//If a key was pressed
	if e.GetType() == sdl.KEYDOWN && (e.(*sdl.KeyboardEvent)).Repeat == 0 {
		//Adjust the velocity
		switch (e.(*sdl.KeyboardEvent)).Keysym.Sym {
		case sdl.K_UP:
			d.mVelY -= DotVel
			break
		case sdl.K_DOWN:
			d.mVelY += DotVel
			break
		case sdl.K_LEFT:
			d.mVelX -= DotVel
			break
		case sdl.K_RIGHT:
			d.mVelX += DotVel
			break
		}
	} else if e.GetType() == sdl.KEYUP && (e.(*sdl.KeyboardEvent)).Repeat == 0 { //If a key was released
		//Adjust the velocity
		switch (e.(*sdl.KeyboardEvent)).Keysym.Sym {
		case sdl.K_UP:
			d.mVelY += DotVel
			break
		case sdl.K_DOWN:
			d.mVelY -= DotVel
			break
		case sdl.K_LEFT:
			d.mVelX += DotVel
			break
		case sdl.K_RIGHT:
			d.mVelX -= DotVel
			break
		}
	}
}

//Move moves the dot and checks collision against tiles
func (d *Dot) Move(tiles []*Tile) {
	//Move the dot left or right
	d.mBox.X += d.mVelX

	//If the dot went too far to the left or right or touched a wall
	if d.mBox.X < 0 || d.mBox.X+DotWidth > levelWidth || touchesWall(d.mBox, tiles) {
		//Move back
		d.mBox.X -= d.mVelX
	}

	//Move the dot up or down
	d.mBox.Y += d.mVelY

	//If the dot went too far up or down or touched a wall
	if d.mBox.Y < 0 || d.mBox.Y+DotHeight > levelHeight || touchesWall(d.mBox, tiles) {
		//Move back
		d.mBox.Y -= d.mVelY
	}
}

//SetCamera centers the camera over the dot
func (d *Dot) SetCamera(camera *sdl.Rect) {
	//Center the camera over the dot
	camera.X = (d.mBox.X + DotWidth/2) - screenWitdh/2
	camera.Y = (d.mBox.Y + DotHeight/2) - screenHeight/2

	//Keep the camera in bounds
	if camera.X < 0 {
		camera.X = 0
	}
	if camera.Y < 0 {
		camera.Y = 0
	}
	if camera.X > levelWidth-camera.W {
		camera.X = levelWidth - camera.W
	}
	if camera.Y > levelHeight-camera.H {
		camera.Y = levelHeight - camera.H
	}
}

//Render shows the dot on the screen
func (d *Dot) Render(camera *sdl.Rect) error {
	//Show the dot
	err := gDotTexture.Render(d.mBox.X-camera.X, d.mBox.Y-camera.Y, nil, 0, nil, sdl.FLIP_NONE)
	if err != nil {
		return fmt.Errorf("could not render dot: %v", err)
	}

	return nil
}
