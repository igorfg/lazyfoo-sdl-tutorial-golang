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
	//The particles
	particles [TotalParticles]*Particle

	//The X and Y offsets of the dot
	mPosX int32
	mPosY int32

	//The velocity of the dot
	mVelX int32
	mVelY int32
}

//NewDot allocates particles
func NewDot() *Dot {
	d := &Dot{}

	//Initialize particles
	for i := 0; i < TotalParticles; i++ {
		d.particles[i] = NewParticle(d.mPosX, d.mPosY)
	}

	return d
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
	if err := gDotTexture.Render(d.mPosX, d.mPosY, nil, 0, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not render dot: %v", err)
	}

	//Show particles on top of the dot
	if err := d.renderParticles(); err != nil {
		return fmt.Errorf("could not render dot's particle: %v", err)
	}

	return nil
}

//Shows the particles
func (d *Dot) renderParticles() error {
	//Go through particles
	for i := 0; i < TotalParticles; i++ {
		//Delete and replace dead particles
		if d.particles[i].IsDead() {
			d.particles[i] = NewParticle(d.mPosX, d.mPosY)
		}
	}

	//Show particles
	for i := 0; i < TotalParticles; i++ {
		if err := d.particles[i].Render(); err != nil {
			return err
		}
	}

	return nil
}
