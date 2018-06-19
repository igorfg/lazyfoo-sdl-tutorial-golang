package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	//Screen dimension constants
	screeWitdh  = 640
	screeHeight = 480
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window
	//The surface contained by the window
	gScreenSurface *sdl.Surface
	//The image we will load and show on the screen
	gXOut *sdl.Surface
)

func main() {
	//Start up SDL and create window
	if err := initSDl(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	//Load media
	if err := loadMedia(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	//Main loop flag
	quit := false

	//Event handler
	var e sdl.Event = nil

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = nil {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}
		//Apply the image
		if err := gXOut.Blit(nil, gScreenSurface, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Could not blit surface! SDL Error: %v\n", err)
		}

		//Update the surface
		if err := gWindow.UpdateSurface(); err != nil {
			fmt.Fprintf(os.Stderr, "Could not update surface! SDL Error: %v\n", err)
		}
	}

	if err := close(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not close SDL! SDL Error: %v\n", err)
	}
}

func initSDl() error {
	//Local error declaration
	var err error

	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("SDL could not initialize! SDL_ERROR: %v", err)
	}
	//Create Window
	gWindow, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screeWitdh, screeHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Window could not be created! SDL_Error: %v", err)
	}
	//Get window surface
	gScreenSurface, err = gWindow.GetSurface()
	if err != nil {
		return fmt.Errorf("Could not get screen surface: %v", err)
	}
	return nil
}

func loadMedia() error {
	//Local error declaration
	var err error

	//Load splash image
	gXOut, err = sdl.LoadBMP("x.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load image %v! SDL Error: %v", "x.bmp", err)
	}
	return nil
}

func close() error {
	//Deallocate surface
	gXOut.Free()
	gXOut = nil

	//Destroy window
	if err := gWindow.Destroy(); err != nil {
		return fmt.Errorf("Coud not destroy window! SDL Error: %v", err)
	}
	gWindow = nil

	//Quit SDL Subsystems
	sdl.Quit()
	return nil
}
