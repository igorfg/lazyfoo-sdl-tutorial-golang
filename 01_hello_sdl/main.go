package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	var (
		//The window we'll be rendering to
		window *sdl.Window
		//The surface contained by the window
		screenSurface *sdl.Surface
	)

	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Errorf("SDL could not initialize! SDL_Error: %v", err)
	} else {
		//Create window
		window, err = sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			screenWidth, screenHeight, sdl.WINDOW_SHOWN)
		if err != nil {
			fmt.Errorf("Window could not be created! SDL_ERROR: %v", err)
		} else {
			//Get window surface
			screenSurface, err = window.GetSurface()
			if err != nil {
				fmt.Errorf("could not get surface! SDL_Error: %v", err)
			}

			//Fill the surface white
			err = screenSurface.FillRect(nil, sdl.MapRGB(screenSurface.Format, 255, 255, 255))
			if err != nil {
				fmt.Errorf("could not fill surface! SDL_Error: %v", err)
			}

			//Update the surface
			err = window.UpdateSurface()
			if err != nil {
				fmt.Errorf("could not update surface! SDL_Error: %v", err)
			}

			//Wait two seconds
			sdl.Delay(2000)
		}
	}
	//Destroy window
	if err := window.Destroy(); err != nil {
		fmt.Errorf("could not destroy window! SDL_Error: %v", err)
	}

	//Quit SDL subsystems
	sdl.Quit()
}
