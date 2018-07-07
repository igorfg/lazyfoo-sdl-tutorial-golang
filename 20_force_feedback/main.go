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

//Analog joystick dead zone
const joystickDeadZone = 8000

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene texture
	gSplashTexture LTexture

	//Game controller 1 handler with force feedback
	gGameController   *sdl.Joystick
	gControllerHaptic *sdl.Haptic
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
			} else if e.GetType() == sdl.JOYBUTTONDOWN { //Joystick buttom press
				//Play rumble at 75% strength for 500 milliseconds
				err := gControllerHaptic.RumblePlay(0.75, 500)
				if err != nil {
					fmt.Printf("Warning: Unable to play rumble! %v\n", err)
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

		err = gSplashTexture.Render(0, 0, nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render splash texture: %v", err)
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
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK | sdl.INIT_HAPTIC); err != nil {
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
		} else {
			//Add controller mapping from https://github.com/gabomdq/SDL_GameControllerDB
			if sdl.GameControllerAddMapping("030000004c050000c405000011010000,PS4 Controller,a:b1,b:b2,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b12,leftshoulder:b4,leftstick:b10,lefttrigger:a3,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:a4,rightx:a2,righty:a5,start:b9,x:b0,y:b3,platform:Linux,") < 0 {
				fmt.Printf("Warning: Controller Mapping could not be added or updated: %v\n", sdl.GetError())
			} else {
				//Get controller haptic device
				gControllerHaptic, err = sdl.HapticOpenFromJoystick(gGameController)
				if err != nil {
					fmt.Printf("Warning: Controller does not support haptics! SDL Error: %v\n", err)
				} else {
					//Get initialize rumble
					err = gControllerHaptic.RumbleInit()
					if err != nil {
						fmt.Printf("Warning: Unable to initialize rumble! SDL Error: %v\n", err)
					}
				}

			}
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
	err = gSplashTexture.LoadFromFile("splash.png")
	if err != nil {
		return fmt.Errorf("failed to load splash texture: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gSplashTexture.Free(); err != nil {
		return fmt.Errorf("could not free splash texture: %v", err)
	}

	//Close game controller
	gControllerHaptic.Close()
	gGameController.Close()
	gGameController = nil
	gControllerHaptic = nil

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
