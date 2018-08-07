package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	//Screen dimension constants
	screenWitdh  = 640
	screenHeight = 480

	//Total Windows
	totalWindows = 3
)

var (
	//Our custom windows
	gWindows [totalWindows]LWindow

	//Globally used font
	gFont *ttf.Font
)

func main() {
	//Start up SDL and create window
	if err := initSDL(); err != nil {
		log.Fatalf("Could not init SDL: %v\n", err)
	}

	for i := 1; i < totalWindows; i++ {
		if err := gWindows[i].Init(); err != nil {
			log.Fatalf("window %d could not be initialized: %v", i, err)
		}
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
			for i := 0; i < totalWindows; i++ {
				if err := gWindows[i].HandleEvent(e); err != nil {
					log.Fatalf("could not handle window %d's event: %v", i, err)
				}
			}

			//Pull up window
			if e.GetType() == sdl.KEYDOWN {
				switch (e.(*sdl.KeyboardEvent)).Keysym.Sym {
				case sdl.K_1:
					gWindows[0].Focus()
					break
				case sdl.K_2:
					gWindows[1].Focus()
					break
				case sdl.K_3:
					gWindows[2].Focus()
					break
				}
			}

			//Update all windows
			for i := 0; i < totalWindows; i++ {
				if err := gWindows[i].Render(); err != nil {
					log.Fatal(err)
				}
			}

			//Check all windows
			allWindowsClosed := true
			for i := 0; i < totalWindows; i++ {
				if gWindows[i].IsShown() {
					allWindowsClosed = false
					break
				}
			}

			//Application closed all windows
			if allWindowsClosed {
				quit = true
			}
		}
	}

	//Free resources and close SDL
	if err := close(); err != nil {
		log.Fatalf("Could not close SDL! SDL Error: %v\n", err)
	}
}

func initSDL() error {
	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("SDL could not initialize! SDL_ERROR: %v", err)
	}

	//Set texture filtering to linear
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Printf("Warning: Linear texture filtering not enabled!")
	}

	//Create Windows
	if err := gWindows[0].Init(); err != nil {
		return fmt.Errorf("Window 0 could not be created: %v", err)
	}

	return nil
}

func close() error {
	//Destroy windows
	for i := 0; i < totalWindows; i++ {
		if err := gWindows[i].Free(); err != nil {
			return fmt.Errorf("could not destroy window %d: %v", i, err)
		}
	}

	sdl.Quit()

	return nil
}
