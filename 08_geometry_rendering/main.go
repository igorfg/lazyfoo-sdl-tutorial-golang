package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
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

	//The window renderer
	gRenderer *sdl.Renderer
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

		//Clear screen
		gRenderer.SetDrawColor(255, 255, 255, 255)
		gRenderer.Clear()

		//Render red filled quad
		fillRect := sdl.Rect{
			X: screenWitdh / 4,
			Y: screenHeight / 4,
			W: screenWitdh / 2,
			H: screenHeight / 2,
		}
		gRenderer.SetDrawColor(255, 0, 0, 255)
		gRenderer.FillRect(&fillRect)

		//Render green outlined quad
		outlineRect := sdl.Rect{
			X: screenWitdh / 6,
			Y: screenHeight / 6,
			W: screenWitdh * 2 / 3,
			H: screenHeight * 2 / 3,
		}
		gRenderer.SetDrawColor(0, 255, 0, 255)
		gRenderer.DrawRect(&outlineRect)

		//Draw blue horizontal line
		gRenderer.SetDrawColor(0, 0, 255, 255)
		gRenderer.DrawLine(0, screenHeight/2, screenWitdh, screenHeight/2)

		//Draw vertical line of yellow dots
		gRenderer.SetDrawColor(255, 255, 0, 255)
		for i := 0; i < screenHeight; i += 4 {
			gRenderer.DrawPoint(screenWitdh/2, int32(i))
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

	//Create renderer for window
	if gRenderer, err = sdl.CreateRenderer(gWindow, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return fmt.Errorf("Renderer could not be created! SDL Error: %v", err)
	}

	//Initialize renderer color
	gRenderer.SetDrawColor(255, 255, 255, 255)

	//Initialize PNG loading
	imgFlags := img.INIT_PNG
	if (img.Init(imgFlags) & imgFlags) == 0 {
		return fmt.Errorf("SDL_image could not initialize! SDL_image Error: %v", img.GetError())
	}

	return nil
}

func loadMedia() error {
	//Nothing to load
	return nil
}

func close() error {
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
	img.Quit()
	sdl.Quit()
	return nil
}

func loadTexture(path string) (*sdl.Texture, error) {
	//The final texture
	var newTexture *sdl.Texture

	//Load image at specified path
	loadedSurface, err := img.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to load image %v! SDL_image: %v", path, err)
	}

	//Create texture from surface pixels
	if newTexture, err = gRenderer.CreateTextureFromSurface(loadedSurface); err != nil {
		return nil, fmt.Errorf("Unable to create texture from %v! SDL Error: %v", path, err)
	}

	//Get rid of old loaded surface
	loadedSurface.Free()

	return newTexture, nil
}
