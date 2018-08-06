package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	//Screen dimension constants
	screenWitdh  = 640
	screenHeight = 480
)

var (
	//Our custom window
	gWindow LWindow

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gSceneTexture LTexture
)

func main() {
	//Start up SDL and create window
	if err := initSDL(); err != nil {
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

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
			//Handle window events
			err := gWindow.HandleEvent(e)
			if err != nil {
				log.Fatal(err)
			}
		}

		//Only draw when not minimized
		if !gWindow.IsMinimized() {
			//Clear screen
			err := gRenderer.SetDrawColor(255, 255, 255, 255)
			if err != nil {
				log.Fatalf("could not set draw color for renderer: %v", err)
			}
			err = gRenderer.Clear()
			if err != nil {
				log.Fatalf("could not clear renderer: %v", err)
			}

			//Render text textures
			err = gSceneTexture.Render((gWindow.MWidth()-gSceneTexture.MWidth())/2,
				(gWindow.MHeight()-gSceneTexture.MHeight())/2, nil, 0, nil, sdl.FLIP_NONE)
			if err != nil {
				log.Fatalf("could not render scene texture: %v\n", err)
			}

			//Update screen
			gRenderer.Present()

		}

	}

	//Free resources and close SDL
	if err := close(); err != nil {
		log.Fatalf("Could not close SDL! SDL Error: %v\n", err)
	}
}

func initSDL() error {
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
	if err := gWindow.Init(); err != nil {
		return fmt.Errorf("Window could not be created! SDL_Error: %v", err)
	}

	//Create vsynced renderer for window
	if gRenderer, err = gWindow.CreateRenderer(); err != nil {
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
	//Local error declaration
	var err error

	if err = gSceneTexture.LoadFromFile("window.png"); err != nil {
		return fmt.Errorf("could not load window texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gSceneTexture.Free(); err != nil {
		return fmt.Errorf("could not free scene texture: %v", err)
	}

	//Destroy window
	if err := gRenderer.Destroy(); err != nil {
		return fmt.Errorf("Could not destroy global renderer: %v", err)
	}
	if err := gWindow.Free(); err != nil {
		return fmt.Errorf("Could not free window: %v", err)
	}
	gRenderer = nil

	//Quit SDL Subsystems
	ttf.Quit()
	img.Quit()
	sdl.Quit()

	return nil
}
