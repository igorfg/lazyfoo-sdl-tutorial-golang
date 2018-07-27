package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//Screen dimension constants
const (
	screenWitdh  = 640
	screenHeight = 480
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gDotTexture LTexture
)

func main() {
	//Start up SDL and create window
	if err := initSDl(); err != nil {
		log.Fatalf("Could not init SDL: %v\n", err)
	}

	//Load media
	if err := loadMedia(); err != nil {
		log.Fatalf("Could not load media: %v\n", err)
	}

	//Main loop flag
	var quit bool

	//Event handler
	var e sdl.Event

	//The dot that will be moving around on the screen
	var dot = NewDot(dotWidth, dotHeight)

	//The dot that will be collided against
	var otherDot = NewDot(screenWitdh/4, screenHeight/4)

	//Set the wall
	wall := &sdl.Rect{X: 300, Y: 40, W: 40, H: 400}

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}

			//Handle input for the dot
			dot.HandleEvent(e)
		}

		//Move the dot and check collision
		dot.Move(wall, otherDot.GetCollider())

		//Clear screen
		err := gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render wall
		err = gRenderer.SetDrawColor(0, 0, 0, 255)
		if err != nil {
			log.Fatalf("coud not set draw color for the wall")
		}
		err = gRenderer.DrawRect(wall)
		if err != nil {
			log.Fatalf("could not draw wall")
		}

		//Render dots
		err = dot.Render()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		err = otherDot.Render()
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		//Update screen
		gRenderer.Present()
	}

	//Free resources and close SDL
	if err := close(); err != nil {
		log.Fatalf("Could not close SDL! SDL Error: %v\n", err)
	}
}

func initSDl() error {
	//Local error declaration
	var err error

	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("SDL could not initialize! SDL_ERROR: %v", err)
	}

	//Set texture filtering to linear
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Printf("Warning: Linear texture filtering not enabled!")
	}

	//Create Window
	gWindow, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWitdh, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Window could not be created! SDL_Error: %v", err)
	}

	//Create vsynced renderer for window
	if gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC); err != nil {
		return fmt.Errorf("Renderer could not be created! SDL Error: %v", err)
	}

	//Initialize renderer color
	gRenderer.SetDrawColor(255, 255, 255, 255)

	//Initialize PNG loading
	imgFlags := img.INIT_PNG
	if (img.Init(imgFlags) & imgFlags) == 0 {
		return fmt.Errorf("SDL_image could not initialize! SDL_image Error: %v", img.GetError())
	}

	//Initialize SDL_ttf
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("SDL_ttf could not initialize! SDL_ttf Error: %v", err)
	}

	return nil
}

func loadMedia() error {
	var err error

	//Load dot texture
	err = gDotTexture.LoadFromFile("dot.bmp")
	if err != nil {
		return fmt.Errorf("Failed to load dot texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gDotTexture.Free(); err != nil {
		return fmt.Errorf("could not free dot texture: %v", err)
	}

	//Destroy window
	if err := gRenderer.Destroy(); err != nil {
		return fmt.Errorf("Could not destroy global renderer: %v", err)
	}
	if err := gWindow.Destroy(); err != nil {
		return fmt.Errorf("Could not destroy window: %v", err)
	}
	gWindow = nil
	gRenderer = nil

	//Quit SDL Subsystems
	ttf.Quit()
	img.Quit()
	sdl.Quit()

	return nil
}

func checkCollision(a *Circle, p Geometry) bool {
	//Checking if collision is with a Circle
	if b, isCircle := p.(*Circle); isCircle {
		//Calculate total radius squared
		totalRadiusSquared := a.R + b.R
		totalRadiusSquared *= totalRadiusSquared

		//If the distance between the centers of the circles is less than the sum of the radii
		if distanceSquared(a.X, a.Y, b.X, b.Y) < float64(totalRadiusSquared) {
			//The circles have collided
			return true
		}
		//If not
		return false
	}

	//Collision is with a SDL_Rect
	b := p.(*sdl.Rect)
	//Closest point of collision box
	var cX, cY int32

	//Find closest x offset
	if a.X < b.X {
		cX = b.X
	} else if a.X > b.X+b.W {
		cX = b.X + b.W
	} else {
		cX = a.X
	}

	//Find closest y offset
	if a.Y < b.Y {
		cY = b.Y
	} else if a.Y > b.Y+b.H {
		cY = b.Y + b.H
	} else {
		cY = a.Y
	}

	//If the closest point is inside the circle
	if distanceSquared(a.X, a.Y, cX, cY) < float64(a.R*a.R) {
		//This box and the circle have collided
		return true
	}
	//If the shapes have not collided
	return false
}

func distanceSquared(x1, y1, x2, y2 int32) float64 {
	deltaX, deltaY := x2-x1, y2-y1
	return float64(deltaX*deltaX + deltaY*deltaY)
}
