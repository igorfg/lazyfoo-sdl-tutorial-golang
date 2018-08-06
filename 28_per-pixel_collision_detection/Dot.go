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

	//Dot's collision boxes
	mColliders []sdl.Rect
}

//NewDot initializes the variables
func NewDot(x, y int32) *Dot {

	//Initialize new Dot with offsets and necessary SDL_rects
	dot := &Dot{
		mPosX:      x,
		mPosY:      y,
		mColliders: make([]sdl.Rect, 11),
	}

	//Initialize the collision boxes' width and height
	dot.mColliders[0].W, dot.mColliders[0].H = 6, 1
	dot.mColliders[1].W, dot.mColliders[1].H = 10, 1
	dot.mColliders[2].W, dot.mColliders[2].H = 14, 1
	dot.mColliders[3].W, dot.mColliders[3].H = 16, 2
	dot.mColliders[4].W, dot.mColliders[4].H = 18, 2
	dot.mColliders[5].W, dot.mColliders[5].H = 20, 6
	dot.mColliders[6].W, dot.mColliders[6].H = 18, 2
	dot.mColliders[7].W, dot.mColliders[7].H = 16, 2
	dot.mColliders[8].W, dot.mColliders[8].H = 14, 1
	dot.mColliders[9].W, dot.mColliders[9].H = 10, 1
	dot.mColliders[10].W, dot.mColliders[10].H = 6, 1

	//Initialize colliders relative to position
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
func (d *Dot) Move(otherColliders []sdl.Rect) {
	//Move the dot left or right
	d.mPosX += d.mVelX
	d.shiftColliders()

	//If the dot collided or went too far to the left or right
	if (d.mPosX < 0) || (d.mPosX+dotWidth > screenWitdh) || checkCollision(d.mColliders, otherColliders) {
		//Move back
		d.mPosX -= d.mVelX
		d.shiftColliders()
	}

	//Move the dot up or down
	d.mPosY += d.mVelY
	d.shiftColliders()

	//If the dot went too far up or down
	if (d.mPosY < 0) || (d.mPosY+dotHeight > screenHeight) || checkCollision(d.mColliders, otherColliders) {
		//Move back
		d.mPosY -= d.mVelY
		d.shiftColliders()
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

//shiftColliders moves the collision boxes relative to the dot's offset
func (d *Dot) shiftColliders() {
	//The row offset
	var r int32

	//Go through the dot's collision boxes
	for set := 0; set < len(d.mColliders); set++ {
		//Center the collision box
		d.mColliders[set].X = d.mPosX + (dotWidth-d.mColliders[set].W)/2

		//Set collision box at its row offset
		d.mColliders[set].Y = d.mPosY + r

		//Move the row offset down the height of the collision box
		r += d.mColliders[set].H
	}

}

//GetColliders gets the collision boxes
func (d Dot) GetColliders() []sdl.Rect {
	return d.mColliders
}
