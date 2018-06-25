package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

//Screen dimension constants
const (
	screenWitdh            = 640
	screenHeight           = 480
	walkingAnimationFrames = 4
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Walking animation
	gSpriteClips        [walkingAnimationFrames]sdl.Rect
	gSpriteSheetTexture LTexture
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

	//Current animtation frame
	var frame int

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
		err := gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render current frame
		currentClip := &gSpriteClips[frame/4]
		err = gSpriteSheetTexture.Render((screenWitdh-currentClip.W)/2,
			(screenHeight-currentClip.H)/2, currentClip)
		if err != nil {
			log.Fatalf("could not render sprite sheet texture: %v", err)
		}

		//Update screen
		gRenderer.Present()

		//Go to next frame
		frame++

		//Cycle animation
		if frame/4 >= walkingAnimationFrames {
			frame = 0
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

	return nil
}

func loadMedia() error {
	//Load sprite sheet texture
	err := gSpriteSheetTexture.LoadFromFile("foo.png")
	if err != nil {
		return fmt.Errorf("failed to load walking animation texture: %v", err)
	}

	//Set sprite clips
	gSpriteClips[0].X = 0
	gSpriteClips[0].Y = 0
	gSpriteClips[0].W = 64
	gSpriteClips[0].H = 205

	gSpriteClips[1].X = 64
	gSpriteClips[1].Y = 0
	gSpriteClips[1].W = 64
	gSpriteClips[1].H = 205

	gSpriteClips[2].X = 128
	gSpriteClips[2].Y = 0
	gSpriteClips[2].W = 64
	gSpriteClips[2].H = 205

	gSpriteClips[3].X = 196
	gSpriteClips[3].Y = 0
	gSpriteClips[3].W = 64
	gSpriteClips[3].H = 205

	return nil
}

func close() error {
	//Free loaded images
	if err := gSpriteSheetTexture.Free(); err != nil {
		return fmt.Errorf("could not free sprite sheet texture: %v", err)
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
