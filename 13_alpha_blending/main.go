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

	//Scene textures
	gModulatedTexture  LTexture
	gBackgroundTexture LTexture
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

	//Modulation component
	var a uint8 = 255

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = nil {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.KEYDOWN { //Handle key presses
				//Increase alpha on w
				if e.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_w {
					//Cap if over 255
					if a+32 < a {
						a = 255
					} else { //Increment otherwise
						a += 32
					}
				} else if e.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_s { //Decrease alpha on s
					//Cap if below -
					if a-32 > a {
						a = 0
					} else { //Decrement otherwise
						a -= 32
					}
				}
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

		//Render background
		err = gBackgroundTexture.Render(0, 0, nil)
		if err != nil {
			log.Fatalf("could not render background texture: %v", err)
		}

		//Render front blended
		err = gModulatedTexture.SetAlpha(a)
		if err != nil {
			log.Fatalf("could not set alpha from texture: %v", err)
		}
		err = gModulatedTexture.Render(0, 0, nil)
		if err != nil {
			log.Fatalf("could not render modulated texture: %v", err)
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
	//Load front alpha texture
	err := gModulatedTexture.LoadFromFile("fadeout.png")
	if err != nil {
		return fmt.Errorf("failed to load front texture: %v", err)
	}

	//Set standard alpha blending
	gModulatedTexture.SetBlendMode(sdl.BLENDMODE_BLEND)

	//Load background texture
	err = gBackgroundTexture.LoadFromFile("fadein.png")
	if err != nil {
		return fmt.Errorf("failed to load background texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gModulatedTexture.Free(); err != nil {
		return fmt.Errorf("could not free modulated texture: %v", err)
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
	newTexture, err = gRenderer.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		return nil, fmt.Errorf("Unable to create texture from %v! SDL Error: %v", path, err)
	}

	//Get rid of old loaded surface
	loadedSurface.Free()

	return newTexture, nil
}
