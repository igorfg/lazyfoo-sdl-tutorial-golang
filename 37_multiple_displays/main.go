package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//Screen dimension constants
	screenWitdh  = 640
	screenHeight = 480
)

var (
	//Our custom window
	gWindow LWindow

	//Display data
	gTotalDisplays int
	gDisplayBounds []sdl.Rect
)

func main() {
	//Start up SDL and create window
	if err := initSDL(); err != nil {
		log.Fatalf("Could not init SDL: %v\n", err)
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
			if err := gWindow.HandleEvent(e); err != nil {
				log.Fatal(err)
			}
		}

		if err := gWindow.Render(); err != nil {
			log.Fatal(err)
		}
	}

	//Free resources and close SDL
	if err := close(); err != nil {
		log.Fatalf("Could not close SDL! SDL Error: %v\n", err)
	}
}

func initSDL() error {
	var err error

	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("SDL could not initialize! SDL_ERROR: %v", err)
	}

	//Set texture filtering to linear
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Printf("Warning: Linear texture filtering not enabled!")
	}

	//Get number of displays
	if gTotalDisplays, err = sdl.GetNumVideoDisplays(); err != nil {
		return fmt.Errorf("could not get number of video displays")
	}
	if gTotalDisplays < 2 {
		fmt.Println("Warning! Only one display connected!")
	}

	//Get bounds of each display
	gDisplayBounds = make([]sdl.Rect, gTotalDisplays)
	for i := 0; i < gTotalDisplays; i++ {
		if gDisplayBounds[i], err = sdl.GetDisplayBounds(i); err != nil {
			return fmt.Errorf("could not get display %d's bounds: %v", i, err)
		}
	}

	//Create Window
	if err := gWindow.Init(); err != nil {
		return fmt.Errorf("Window could not be created: %v", err)
	}

	return nil
}

func close() error {
	//Destroy window
	if err := gWindow.Free(); err != nil {
		return fmt.Errorf("could not destroy window: %v", err)
	}

	sdl.Quit()

	return nil
}
