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
	gPressTexture LTexture
	gUpTexture    LTexture
	gDownTexture  LTexture
	gLeftTexture  LTexture
	gRightTexture LTexture
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

	//Current rendered texture
	var currentTexture *LTexture

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = nil {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}

		//Set texture based on current keystate
		currentKeyStates := sdl.GetKeyboardState()
		if currentKeyStates[sdl.SCANCODE_UP] != 0 {
			currentTexture = &gUpTexture
		} else if currentKeyStates[sdl.SCANCODE_DOWN] != 0 {
			currentTexture = &gDownTexture
		} else if currentKeyStates[sdl.SCANCODE_LEFT] != 0 {
			currentTexture = &gLeftTexture
		} else if currentKeyStates[sdl.SCANCODE_RIGHT] != 0 {
			currentTexture = &gRightTexture
		} else {
			currentTexture = &gPressTexture
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

		//Render texture
		err = currentTexture.Render(0, 0, nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render current texture: %v", err)
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

	return nil
}

func loadMedia() error {
	//Local error declaration
	var err error

	//Load press texture
	err = gPressTexture.LoadFromFile("press.png")
	if err != nil {
		return fmt.Errorf("failed to load press texture: %v", err)
	}

	//Load up texture
	err = gUpTexture.LoadFromFile("up.png")
	if err != nil {
		return fmt.Errorf("failed to load up texture: %v", err)
	}

	//Load down texture
	err = gDownTexture.LoadFromFile("down.png")
	if err != nil {
		return fmt.Errorf("failed to load down texture: %v", err)
	}

	err = gLeftTexture.LoadFromFile("left.png")
	if err != nil {
		return fmt.Errorf("failed to load left texture: %v", err)
	}

	err = gRightTexture.LoadFromFile("right.png")
	if err != nil {
		return fmt.Errorf("failed to load right texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gPressTexture.Free(); err != nil {
		return fmt.Errorf("could not free press texture: %v", err)
	}
	if err := gUpTexture.Free(); err != nil {
		return fmt.Errorf("could not free up texture: %v", err)
	}
	if err := gDownTexture.Free(); err != nil {
		return fmt.Errorf("could not free down texture: %v", err)
	}
	if err := gLeftTexture.Free(); err != nil {
		return fmt.Errorf("could not free left texture: %v", err)
	}
	if err := gRightTexture.Free(); err != nil {
		return fmt.Errorf("could not free right texture: %v", err)
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
