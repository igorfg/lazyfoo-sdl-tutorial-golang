package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//dotWidth is the dot's width
	dotWidth = 20
	//dotHeight is the dot's height
	dotHeight = 20

	//DotVel is the maximum axis velocity of the dot
	DotVel = 1
)

//Dot is the dot that will move around on the screen
type Dot struct {
	//The X and Y offsets of the dot
	mPosX, mPosY int32

	//The velocity of the dot
	mVelX, mVelY int32

	//Dot's collision circle
	mCollider Circle
}

//NewDot initializes the variables
func NewDot(x, y int32) *Dot {
	//Initialize new Dot with offsets and necessary SDL_rects
	dot := &Dot{
		mPosX:     x,
		mPosY:     y,
		mCollider: Circle{R: dotWidth / 2},
	}

	//Move collider relative to the circle
	dot.shiftColliders()

	return dot
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
func (d *Dot) Move(square *sdl.Rect, circle *Circle) {
	//Move the dot left or right
	d.mPosX += d.mVelX
	d.shiftColliders()

	//If the dot collided or went too far to the left or right
	if (d.mPosX-d.mCollider.R < 0) || (d.mPosX+d.mCollider.R > screenWitdh) ||
		checkCollision(&d.mCollider, square) || checkCollision(&d.mCollider, circle) {
		//Move back
		d.mPosX -= d.mVelX
		d.shiftColliders()
	}

	//Move the dot up or down
	d.mPosY += d.mVelY
	d.shiftColliders()

	//If the dot went too far up or down
	if (d.mPosY-d.mCollider.R < 0) || (d.mPosY+d.mCollider.R > screenHeight) ||
		checkCollision(&d.mCollider, square) || checkCollision(&d.mCollider, circle) {
		//Move back
		d.mPosY -= d.mVelY
		d.shiftColliders()
	}
}

//Render shows the dot on the screen
func (d *Dot) Render() error {
	//Show the dot
	err := gDotTexture.Render(d.mPosX-d.mCollider.R, d.mPosY-d.mCollider.R, nil, 0, nil, sdl.FLIP_NONE)
	if err != nil {
		return fmt.Errorf("could not render dot: %v", err)
	}

	return nil
}

//shiftColliders moves the collision boxescircle relative to the dot's offset
func (d *Dot) shiftColliders() {
	//Alling the collider with the center of the dot
	d.mCollider.X = d.mPosX
	d.mCollider.Y = d.mPosY
}

//GetCollider gets the collision circle
func (d Dot) GetCollider() *Circle {
	return &d.mCollider
}
