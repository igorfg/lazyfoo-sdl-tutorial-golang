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
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The surface contained by the window
	gScreenSurface *sdl.Surface

	//Current displayed image
	gStretchedSurface *sdl.Surface
)

func main() {
	//Start up SDL and create window
	if err := initSDl(); err != nil {
		log.Fatalf("%v\n", err)
	}

	//Load media
	if err := loadMedia(); err != nil {
		log.Fatalf("%v\n", err)
	}

	//Main loop flag
	quit := false

	//Event handler
	var e sdl.Event

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = nil {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}

		//Apply the image stretched
		var stretchRect = sdl.Rect{X: 0, Y: 0, W: screenWitdh, H: screenHeight}
		if err := gStretchedSurface.BlitScaled(nil, gScreenSurface, &stretchRect); err != nil {
			log.Fatalf("Could not blit surface! SDL Error: %v\n", err)
		}

		//Update the surface
		if err := gWindow.UpdateSurface(); err != nil {
			log.Fatalf("Could not update surface! SDL Error: %v\n", err)
		}
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
	//Create Window
	gWindow, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWitdh, screenHeight, sdl.WINDOW_SHOWN)
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

	//Load stretching surface
	gStretchedSurface, err = loadSurface("stretch.bmp")
	if err != nil {
		return fmt.Errorf("Failed to load stretching image! SDL_ERROR: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded image
	gStretchedSurface.Free()
	gScreenSurface = nil

	//Destroy window
	if err := gWindow.Destroy(); err != nil {
		return fmt.Errorf("Coud not destroy window! SDL Error: %v", err)
	}
	gWindow = nil

	//Quit SDL Subsystems
	sdl.Quit()
	return nil
}

func loadSurface(path string) (*sdl.Surface, error) {
	//The final optimized image
	var optimizedSurface *sdl.Surface

	//Load image at specified path
	loadedSurface, err := sdl.LoadBMP(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to load image %v! SDL_ERROR: %v", path, err)
	}

	//Convert surface to screen format
	optimizedSurface, err = loadedSurface.Convert(gScreenSurface.Format, 0)
	if err != nil {
		return nil, fmt.Errorf("Unable to optimize image %v! SDL Error: %v", path, err)
	}

	//Get rid of old loaded surface
	loadedSurface.Free()

	return optimizedSurface, nil
}
