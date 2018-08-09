package main

import (
	"encoding/binary"
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
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gFooTexture LTexture
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

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
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

		//Render dot
		err = gFooTexture.Render((screenWitdh-gFooTexture.mWidth)/2, (screenHeight-gFooTexture.MHeight())/2,
			nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render foo texture: %v\n", err)
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

	//Load foo' texture
	if err = gFooTexture.LoadFromFile("foo.png"); err != nil {
		return fmt.Errorf("Failed to load foo texture: %v", err)
	}

	//Lock texture
	if err = gFooTexture.LockTexture(); err != nil {
		return fmt.Errorf("unable to lock foo texture: %v", err)
	}

	//Allocate format from window
	format, err := gWindow.GetPixelFormat()
	if err != nil {
		return fmt.Errorf("could not get pixel format from window: %v", err)
	}

	mappingFormat, err := sdl.AllocFormat(uint(format))
	if err != nil {
		return fmt.Errorf("could not get mapping format: %v", err)
	}

	//Get pixel data
	bytePixels := gFooTexture.MPixels()
	pixels := make([]uint32, len(bytePixels)/4)
	for i := range pixels {
		//Assuming little endian
		pixels[i] = uint32(binary.LittleEndian.Uint32(bytePixels[i*4 : (i+1)*4]))
	}

	pixelCount := (gFooTexture.MPitch() / 4) * int(gFooTexture.MHeight())

	//Map colors
	colorKey := sdl.MapRGB(mappingFormat, 0, 255, 255)
	transparent := sdl.MapRGBA(mappingFormat, 255, 255, 255, 0)

	//Color key pixels
	for i := 0; i < pixelCount; i++ {
		if pixels[i] == colorKey {
			binary.LittleEndian.PutUint32(bytePixels[i*4:(i*4)+4], transparent)
		}
	}

	//Unlock texture
	if err = gFooTexture.UnlockTexture(); err != nil {
		return fmt.Errorf("could not unlock foo texture: %v", err)
	}

	//Free format
	mappingFormat.Free()

	return nil
}

func close() error {
	//Free loaded images
	if err := gFooTexture.Free(); err != nil {
		return fmt.Errorf("could not free foo texture: %v", err)
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
