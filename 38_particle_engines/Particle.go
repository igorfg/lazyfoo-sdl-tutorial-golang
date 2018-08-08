package main

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

//TotalParticles is the particle count
const TotalParticles = 20

//Particle is used to make a little animation that follows the Dot around
type Particle struct {
	//Offset
	mPosX int32
	mPosY int32

	//Current frame of animation
	mFrame int

	//Type of particle
	mTexture *LTexture
}

//NewParticle initializes position and animation
func NewParticle(x, y int32) *Particle {
	p := &Particle{}

	//Set offsets
	p.mPosX = x - 5 + rand.Int31n(25)
	p.mPosY = y - 5 + rand.Int31n(25)

	//Initialize animation
	p.mFrame = rand.Intn(5)

	//Set type
	switch rand.Intn(3) {
	case 0:
		p.mTexture = &gRedTexture
		break
	case 1:
		p.mTexture = &gGreenTexture
		break
	case 2:
		p.mTexture = &gBlueTexture
		break
	}

	return p
}

//Render shows the particle
func (p *Particle) Render() error {
	//Show image
	if err := p.mTexture.Render(p.mPosX, p.mPosY, nil, 0, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not show particle image: %v", err)
	}

	//Show shimmer
	if p.mFrame%2 == 0 {
		if err := gShimmerTexture.Render(p.mPosX, p.mPosY, nil, 0, nil, sdl.FLIP_NONE); err != nil {
			return fmt.Errorf("could not show particle shimmer: %v", err)
		}
	}

	//Animate
	p.mFrame++

	return nil
}

//IsDead checks if particle is dead
func (p *Particle) IsDead() bool {
	return p.mFrame > 10
}
