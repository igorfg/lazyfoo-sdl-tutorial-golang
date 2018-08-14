package main

import (
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
	gTargetTexture LTexture
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

	//Rotation variablles
	var angle float64
	screenCenter := sdl.Point{X: screenWitdh / 2, Y: screenHeight / 2}

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}

		//Rotate
		angle += 2
		if angle > 360 {
			angle -= 360
		}

		//Set self as render target
		if err := gTargetTexture.SetAsRenderTarget(); err != nil {
			log.Fatal(err)
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

		//Render red filled quad
		fillRect := sdl.Rect{X: screenWitdh / 4, Y: screenHeight / 4, W: screenWitdh / 2, H: screenHeight / 2}
		if err := gRenderer.SetDrawColor(255, 0, 0, 0); err != nil {
			log.Fatalf("could not set draw color to red: %v\n", err)
		}
		if err := gRenderer.FillRect(&fillRect); err != nil {
			log.Fatalf("could not fill red rect: %v\n", err)
		}

		//Render green outlined quad
		outlineRect := sdl.Rect{X: screenWitdh / 6, Y: screenHeight / 6, W: screenWitdh * 2 / 3, H: screenHeight * 2 / 3}
		if err := gRenderer.SetDrawColor(0, 255, 0, 255); err != nil {
			log.Fatalf("could not set draw color to green: %v\n", err)
		}
		if err := gRenderer.DrawRect(&outlineRect); err != nil {
			log.Fatalf("could not draw green rect: %v\n", err)
		}

		//Draw blue horizontal line
		if err := gRenderer.SetDrawColor(0, 0, 255, 255); err != nil {
			log.Fatalf("could not set draw color to blue: %v\n", err)
		}
		if err := gRenderer.DrawLine(0, screenHeight/2, screenWitdh, screenHeight/2); err != nil {
			log.Fatalf("could not draw blue line: %v\n", err)
		}

		//Draw vertical line of yellow dots
		if err := gRenderer.SetDrawColor(255, 255, 0, 255); err != nil {
			log.Fatalf("coud not set draw color to yellow: %v\n", err)
		}
		for i := 0; i < screenHeight; i += 4 {
			if err := gRenderer.DrawPoint(screenWitdh/2, int32(i)); err != nil {
				log.Fatalf("could not draw point number %d: %v\n", i, err)
			}
		}

		//Reset render target
		if err := gRenderer.SetRenderTarget(nil); err != nil {
			log.Fatalf("could not reset render target: %v\n", err)
		}

		//Show rendered to texture
		if err = gTargetTexture.Render(0, 0, nil, angle, &screenCenter, sdl.FLIP_NONE); err != nil {
			log.Fatal(err)
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
	//Load blank texture
	if err := gTargetTexture.CreateBlank(screenWitdh, screenHeight, sdl.TEXTUREACCESS_TARGET); err != nil {
		return fmt.Errorf("failed to create target texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gTargetTexture.Free(); err != nil {
		return fmt.Errorf("could not free target texture: %v", err)
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
