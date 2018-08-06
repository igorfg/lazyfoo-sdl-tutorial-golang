package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//The dimensions of the level
const (
	LevelWidth  = 1280
	LevelHeight = 960
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
	gBGTexture  LTexture
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
	var dot Dot

	//The camera area
	camera := sdl.Rect{X: 0, Y: 0, W: screenWitdh, H: screenHeight}

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
		dot.Move()

		//Center the camera over the dot
		camera.X = (dot.GetPosX() + DotWidth/2) - screenWitdh/2
		camera.Y = (dot.GetPosY() + DotHeight/2) - screenHeight/2

		//Keep the camera in bounds
		if camera.X < 0 {
			camera.X = 0
		}
		if camera.Y < 0 {
			camera.Y = 0
		}
		if camera.X > LevelWidth-camera.W {
			camera.X = LevelWidth - camera.W
		}
		if camera.Y > LevelHeight-camera.H {
			camera.Y = LevelHeight - camera.H
		}

		//Clear screen
		err := gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render background
		err = gBGTexture.Render(0, 0, &camera, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render background texture: %v\n", err)
		}

		//Render dot
		err = dot.Render(camera.X, camera.Y)
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
	//Load background texture
	err := gBGTexture.LoadFromFile("bg.png")
	if err != nil {
		return fmt.Errorf("Failed to load background texute: %v", err)
	}

	//Load dot texture
	err = gDotTexture.LoadFromFile("dot.bmp")
	if err != nil {
		return fmt.Errorf("Failed to load dot texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gBGTexture.Free(); err != nil {
		return fmt.Errorf("could not free backgraound texture: %v", err)
	}
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
