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
	var dot = NewDot()

	//Set the wall
	wall := sdl.Rect{X: 300, Y: 40, W: 40, H: 400}

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

		//Move the dot
		dot.Move(&wall)

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
			log.Fatalf("could not draw color for wall rendering: %vzn", err)
		}
		err = gRenderer.DrawRect(&wall)
		if err != nil {
			log.Fatalf("could not render rect: %v\n", err)
		}

		//Render dot
		err = dot.Render()
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

func checkCollision(a, b sdl.Rect) bool {
	//The sides of the rectangles
	var (
		leftA, leftB     int32
		rightA, rightB   int32
		topA, topB       int32
		bottomA, bottomB int32
	)

	//Calculate the sides of rect A
	leftA = a.X
	rightA = a.X + a.W
	topA = a.Y
	bottomA = a.Y + a.H

	//Calculate the sides of rect B
	leftB = b.X
	rightB = b.X + b.W
	topB = b.Y
	bottomB = b.Y + b.H

	//If any of the sides from A are outside of B
	if bottomA <= topB {
		return false
	}
	if topA >= bottomB {
		return false
	}
	if rightA <= leftB {
		return false
	}
	if leftA >= rightB {
		return false
	}

	//If none of the sides from A are outside of B
	return true
}
