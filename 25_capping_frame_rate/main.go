package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//Screen dimension constants
const (
	screenWitdh         = 640
	screenHeight        = 480
	screenFPS           = 60
	screenTicksPerFrame = 1000 / screenFPS
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gFPSTextTexture LTexture
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

	//Set text coloar as black
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//The frames per second timer
	var fpsTimer LTimer

	//The frames per second cap timer
	var capTimer LTimer

	//In memory text stream
	var timeText *bytes.Buffer

	//Start counting frames per second
	var countedFrames int
	fpsTimer.Start()

	//While application is running
	for !quit {
		//Start cap timer
		capTimer.Start()

		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}

		//Calculate and correct FPS
		avgFPS := float64(countedFrames) / (float64(fpsTimer.GetTicks()) / 1000)
		if avgFPS > 2000000 {
			avgFPS = 0
		}

		//Set text to be rendered
		timeText = bytes.NewBufferString("")
		fmt.Fprintf(timeText, "Average Frames Per Second (With Cap) %.4f", avgFPS)

		//Render text
		err := gFPSTextTexture.loadFromRenderedText(timeText.String(), textColor)
		if err != nil {
			log.Fatalf("Unable to render FPS texture: %v\n", err)
		}

		//Clear screen
		err = gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render textures
		err = gFPSTextTexture.Render(
			(screenWitdh-gFPSTextTexture.GetWidth())/2, (screenHeight-gFPSTextTexture.GetHeight())/2,
			nil, 0, nil, sdl.FLIP_NONE,
		)
		if err != nil {
			log.Fatalf("could not render FPS texture: %v\n", err)
		}

		//Update screen
		gRenderer.Present()
		countedFrames++

		//If frame finished early
		frameTicks := capTimer.GetTicks()
		if frameTicks < screenTicksPerFrame {
			//Wait remaining time
			sdl.Delay(screenTicksPerFrame - frameTicks)
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

	//Initialize SDL_ttf
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("SDL_ttf could not initialize! SDL_ttf Error: %v", err)
	}
	return nil
}

func loadMedia() error {
	//Local error declaration
	var err error

	//Open the font
	gFont, err = ttf.OpenFont("lazy.ttf", 28)
	if err != nil {
		return fmt.Errorf("Failed to load lazy font! SDL_ttf Error: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gFPSTextTexture.Free(); err != nil {
		return fmt.Errorf("could not free FPS texture: %v", err)
	}

	//Free global font
	gFont.Close()
	gFont = nil

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
