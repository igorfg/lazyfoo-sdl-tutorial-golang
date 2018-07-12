package main

import (
	"fmt"
	"log"
	"math"

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

	//Arrow texture
	gArrowTexture LTexture

	//Game controller 1 handler
	gGameController *sdl.Joystick
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

	//Normalized direction
	var xDir int
	var yDir int

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.JOYAXISMOTION {
				//Motion on controller 0
				if e.(*sdl.JoyAxisEvent).Which == 0 {
					//X axis motion
					if e.(*sdl.JoyAxisEvent).Axis == 0 {
						//Left of dead zone
						if e.(*sdl.JoyAxisEvent).Value < -joystickDeadZone {
							xDir = -1
						} else if e.(*sdl.JoyAxisEvent).Value > joystickDeadZone { //Right of dead zone
							xDir = 1
						} else {
							xDir = 0
						}
					} else if e.(*sdl.JoyAxisEvent).Axis == 1 { //Y axis motion
						//Below of dead zone
						if e.(*sdl.JoyAxisEvent).Value < -joystickDeadZone {
							yDir = -1
						} else if e.(*sdl.JoyAxisEvent).Value > joystickDeadZone { //Above of dead zone
							yDir = 1
						} else {
							yDir = 0
						}
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

		//Calculate angle
		var joyStickAngle = math.Atan2(float64(yDir), float64(xDir)) * (180.0 / math.Pi)

		//Correct angle
		if xDir == 0 && yDir == 0 {
			joyStickAngle = 0
		}

		//Render joystick 8 way angle
		err = gArrowTexture.Render((screenWitdh-gArrowTexture.GetWidth())/2,
			(screenHeight-gArrowTexture.GetHeight())/2, nil, joyStickAngle, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render arrow texture on screen: %v", err)
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
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK); err != nil {
		return fmt.Errorf("SDL could not initialize! SDL_ERROR: %v", err)
	}

	//Set texture filtering to linear
	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1") {
		fmt.Printf("Warning: Linear texture filtering not enabled!")
	}

	//Check for joysticks
	if sdl.NumJoysticks() < 1 {
		fmt.Printf("Warning: No joysticks connected!\n")
	} else {
		//Load joystick
		gGameController = sdl.JoystickOpen(0)
		if gGameController == nil {
			fmt.Printf("Waning: Unable to open game controller! SDL_Error: %v\n", sdl.GetError())
		}
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
	err = gArrowTexture.LoadFromFile("arrow.png")
	if err != nil {
		return fmt.Errorf("failed to load arrow texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gArrowTexture.Free(); err != nil {
		return fmt.Errorf("could not free arrow texture: %v", err)
	}

	//Close game controller
	gGameController.Close()
	gGameController = nil

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
