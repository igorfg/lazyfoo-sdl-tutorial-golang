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
	DotVel = 20
)

//Dot is the dot that will move around on the screen
type Dot struct {
	mPosX, mPosY int32
	mVelX, mVelY int32
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

//Move moves the dot
func (d *Dot) Move() {
	//Move the dot left or right
	d.mPosX += d.mVelX

	//If the dot went too far to the left or right
	if d.mPosX < 0 || d.mPosX+DotWidth > screenWitdh {
		//Move back
		d.mPosX -= d.mVelX
	}

	//Move the dot up or down
	d.mPosY += d.mVelY

	//If the dot went too far up or down
	if d.mPosY < 0 || d.mPosY+DotHeight > screenHeight {
		//Move back
		d.mPosY -= d.mVelY
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
