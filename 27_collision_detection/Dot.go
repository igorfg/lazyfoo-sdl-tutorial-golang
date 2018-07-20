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
	//The X and Y offsets of the dot
	mPosX, mPosY int32

	//The velocity of the dot
	mVelX, mVelY int32

	//Dot's collision box
	mCollider sdl.Rect
}

//NewDot initializes the variables
func NewDot() *Dot {
	//Set collision box dimension
	collider := sdl.Rect{W: DotWidth, H: DotHeight}

	return &Dot{mCollider: collider}
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

//Move moves the dot and checks collision
func (d *Dot) Move(wall *sdl.Rect) {
	//Move the dot left or right
	d.mPosX += d.mVelX
	d.mCollider.X = d.mPosX

	//If the dot collided or went too far to the left or right
	if (d.mPosX < 0) || (d.mPosX+DotWidth > screenWitdh) || checkCollision(d.mCollider, *wall) {
		//Move back
		d.mPosX -= d.mVelX
		d.mCollider.X = d.mPosX
	}

	//Move the dot up or down
	d.mPosY += d.mVelY
	d.mCollider.Y = d.mPosY

	//If the dot went too far up or down
	if (d.mPosY < 0) || (d.mPosY+DotHeight > screenHeight) || checkCollision(d.mCollider, *wall) {
		//Move back
		d.mPosY -= d.mVelY
		d.mCollider.Y = d.mPosY
	}
}

//Render shows the dot on the screen
func (d *Dot) Render() error {
	//Show the dot
	err := gDotTexture.Render(d.mPosX, d.mPosY, nil, 0, nil, sdl.FLIP_NONE)
	if err != nil {
		return fmt.Errorf("could not render dot: %v", err)
	}

	return nil
}
