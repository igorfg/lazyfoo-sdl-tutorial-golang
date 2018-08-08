package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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
	gDotTexture     LTexture
	gRedTexture     LTexture
	gGreenTexture   LTexture
	gBlueTexture    LTexture
	gShimmerTexture LTexture
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
	var dot = NewDot()

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

		//Move the dot
		dot.Move()

		//Clear screen
		err := gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render objects
		err = dot.Render()
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
	var err error

	//Load dot texture
	if err = gDotTexture.LoadFromFile("dot.bmp"); err != nil {
		return fmt.Errorf("Failed to load dot texture: %v", err)
	}

	//Load red texture
	if err = gRedTexture.LoadFromFile("red.bmp"); err != nil {
		return fmt.Errorf("Failed to load red texture: %v", err)
	}

	//Load green texture
	if err = gGreenTexture.LoadFromFile("green.bmp"); err != nil {
		return fmt.Errorf("Failed to load green texture: %v", err)
	}

	//Load blue texture
	if err = gBlueTexture.LoadFromFile("blue.bmp"); err != nil {
		return fmt.Errorf("Failed to load blue texture: %v", err)
	}

	//Load shimmer texture
	if err = gShimmerTexture.LoadFromFile("shimmer.bmp"); err != nil {
		return fmt.Errorf("Failed to load shimmer texture: %v", err)
	}

	//Set texture transparency
	if err = gRedTexture.SetAlpha(192); err != nil {
		return fmt.Errorf("could not set red texture's alpha")
	}
	if err = gGreenTexture.SetAlpha(192); err != nil {
		return fmt.Errorf("could not set green texture's alpha")
	}
	if err = gBlueTexture.SetAlpha(192); err != nil {
		return fmt.Errorf("could not set blue texture's alpha")
	}
	if err = gShimmerTexture.SetAlpha(192); err != nil {
		return fmt.Errorf("could not set shimmer texture's alpha")
	}

	return nil
}

func close() error {
	//Free loaded images
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
