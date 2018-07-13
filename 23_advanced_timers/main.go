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
	screenWitdh  = 640
	screenHeight = 480
)

//Analog joystick dead zone
const joystickDeadZone = 8000

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gTimeTexture        LTexture
	gPausePromptTexture LTexture
	gStartPromptTexture LTexture
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

	//The application timer
	var timer LTimer

	//In memory text stream
	var timeText *bytes.Buffer

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.KEYDOWN {
				//Start / stop
				if e.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_s {
					if timer.IsStarted() {
						timer.Stop()
					} else {
						timer.Start()
					}
				} else if e.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_p { //Pause / unpause
					if timer.IsPaused() {
						timer.Unpause()
					} else {
						timer.Pause()
					}
				}
			}
		}

		//Set text to be rendered
		timeText = bytes.NewBufferString("")
		fmt.Fprint(timeText, "Seconds since start time ", float64(timer.GetTicks())/float64(1000))

		//Render text
		err := gTimeTexture.loadFromRenderedText(timeText.String(), textColor)
		if err != nil {
			log.Fatalf("Unable to render time texture: %v\n", err)
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

		err = gStartPromptTexture.Render(
			(screenWitdh-gStartPromptTexture.GetWidth())/2, 0, nil, 0, nil, sdl.FLIP_NONE,
		)
		if err != nil {
			log.Fatalf("could not render start prompt texture: %v", err)
		}

		err = gPausePromptTexture.Render(
			(screenWitdh-gPausePromptTexture.GetWidth())/2, gStartPromptTexture.GetHeight(),
			nil, 0, nil, sdl.FLIP_NONE,
		)
		if err != nil {
			log.Fatalf("could not render start prompt texture: %v", err)
		}

		err = gTimeTexture.Render((screenWitdh-gTimeTexture.GetWidth())/2,
			(screenHeight-gTimeTexture.GetHeight())/2, nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render time texture: %v", err)
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
	//Local error declaration
	var err error

	//Open the font
	gFont, err = ttf.OpenFont("lazy.ttf", 28)
	if err != nil {
		return fmt.Errorf("Failed to load lazy font! SDL_ttf Error: %v", err)
	}

	//Set text color as black
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//Load prompt texture
	err = gStartPromptTexture.loadFromRenderedText("Press S to Start or Stop the Timer", textColor)
	if err != nil {
		return fmt.Errorf("Unable to render start/stop prompt texture: %v", err)
	}

	err = gPausePromptTexture.loadFromRenderedText("Press P to Pause or Unpause the Timer", textColor)
	if err != nil {
		return fmt.Errorf("Unable to render pause/unpause prompt texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gStartPromptTexture.Free(); err != nil {
		return fmt.Errorf("could not free start/stop prompt texture: %v", err)
	}
	if err := gPausePromptTexture.Free(); err != nil {
		return fmt.Errorf("could not free pause/unpause prompt texture: %v", err)
	}
	if err := gTimeTexture.Free(); err != nil {
		return fmt.Errorf("could not free prompt texture: %v", err)
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
