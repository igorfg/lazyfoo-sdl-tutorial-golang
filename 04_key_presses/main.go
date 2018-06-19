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

	//Key press surfaces constants
	keyPressSurfaceDefault = iota
	keyPressSurfaceUp
	keyPressSurfaceDown
	keyPressSurfaceLeft
	keyPressSurfaceRight
	keyPressSurfaceTotal
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The surface contained by the window
	gScreenSurface *sdl.Surface

	//The image we will load and show on the screen
	gKeyPressSurfaces [keyPressSurfaceTotal]*sdl.Surface

	//Current displayed image
	gCurrentSurface *sdl.Surface
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
	var e sdl.Event

	//Set default current surface
	gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceDefault]

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = nil {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.KEYDOWN { //User presses key
				//Select surfaces based on key press
				switch e.(*sdl.KeyboardEvent).Keysym.Sym {
				case sdl.K_UP:
					gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceUp]
					break
				case sdl.K_DOWN:
					gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceDown]
					break
				case sdl.K_LEFT:
					gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceLeft]
					break
				case sdl.K_RIGHT:
					gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceRight]
					break
				default:
					gCurrentSurface = gKeyPressSurfaces[keyPressSurfaceDefault]
					break
				}
			}
		}
		//Apply the image
		if err := gCurrentSurface.Blit(nil, gScreenSurface, nil); err != nil {
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

	//Load default surface
	gKeyPressSurfaces[keyPressSurfaceDefault], err = loadSurface("press.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load surface %v! SDL Error: %v", "press.bmp", err)
	}

	//Load up surface
	gKeyPressSurfaces[keyPressSurfaceUp], err = loadSurface("up.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load surface %v! SDL Error: %v", "up.bmp", err)
	}

	//Load down surface
	gKeyPressSurfaces[keyPressSurfaceDown], err = loadSurface("down.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load surface %v! SDL Error: %v", "down.bmp", err)
	}

	//Load left surface
	gKeyPressSurfaces[keyPressSurfaceLeft], err = loadSurface("left.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load surface %v! SDL Error: %v", "left.bmp", err)
	}

	//Load right surface
	gKeyPressSurfaces[keyPressSurfaceRight], err = loadSurface("right.bmp")
	if err != nil {
		return fmt.Errorf("Unable to load surface %v! SDL Error: %v", "right.bmp", err)
	}

	return nil
}

func close() error {
	//Deallocate surfaces
	for i := 0; i < keyPressSurfaceTotal; i++ {
		gKeyPressSurfaces[i].Free()
		gKeyPressSurfaces[i] = nil
	}

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
	//Load image at specified path

	loadedSurface, err := sdl.LoadBMP(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to load image %v! SDL Error: %v", path, err)
	}
	return loadedSurface, nil
}
