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

	//Scene sprites
	gSpriteClips        [4]sdl.Rect
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

		//Render top left sprite
		gSpriteSheetTexture.Render(0, 0, &gSpriteClips[0])

		//Render top right sprite
		gSpriteSheetTexture.Render(screenWitdh-gSpriteClips[1].W, 0, &gSpriteClips[1])

		//Render bottom left sprite
		gSpriteSheetTexture.Render(0, screenHeight-gSpriteClips[2].H, &gSpriteClips[2])

		//Render bottom right sprite
		gSpriteSheetTexture.Render(screenWitdh-gSpriteClips[3].W,
			screenHeight-gSpriteClips[3].H, &gSpriteClips[3])

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
	//Load Foo' texture
	err := gSpriteSheetTexture.LoadFromFile("dots.png")
	if err != nil {
		return fmt.Errorf("failed to load sprite sheet texture: %v", err)
	}

	//Set top left sprite
	gSpriteClips[0].X = 0
	gSpriteClips[0].Y = 0
	gSpriteClips[0].W = 100
	gSpriteClips[0].H = 100

	//Set top right sprite
	gSpriteClips[1].X = 100
	gSpriteClips[1].Y = 0
	gSpriteClips[1].W = 100
	gSpriteClips[1].H = 100

	//Set bottom left sprite
	gSpriteClips[2].X = 0
	gSpriteClips[2].Y = 100
	gSpriteClips[2].W = 100
	gSpriteClips[2].H = 100

	//Set bottom right sprite
	gSpriteClips[3].X = 100
	gSpriteClips[3].Y = 100
	gSpriteClips[3].W = 100
	gSpriteClips[3].H = 100

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
