package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/mix"
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
	gPromptTexture LTexture

	//The music that will be played
	gMusic *mix.Music

	//The sound effects that will be used
	gScratch *mix.Chunk
	gHigh    *mix.Chunk
	gMedium  *mix.Chunk
	gLow     *mix.Chunk
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
			} else if e.GetType() == sdl.KEYDOWN { //Handle key press
				switch e.(*sdl.KeyboardEvent).Keysym.Sym {
				//Play high sound effect
				case sdl.K_1:
					_, err := gHigh.Play(-1, 0)
					if err != nil {
						fmt.Printf("Warning! Could not play high sound effect: %v\n", err)
					}
					break
				//Play medium sound effect
				case sdl.K_2:
					_, err := gMedium.Play(-1, 0)
					if err != nil {
						fmt.Printf("Warning! Could not play medium sound effect: %v\n", err)
					}
					break
				//Play low sound effect
				case sdl.K_3:
					_, err := gLow.Play(-1, 0)
					if err != nil {
						fmt.Printf("Warning! Could not play low sound effect: %v\n", err)
					}
					break
				//Play scratch sound effect
				case sdl.K_4:
					_, err := gScratch.Play(-1, 0)
					if err != nil {
						fmt.Printf("Warning! Could not play scratch sound effect: %v\n", err)
					}
					break
				case sdl.K_9:
					//If there is no music playing
					if !mix.PlayingMusic() {
						//Play the music
						err := gMusic.Play(-1)
						if err != nil {
							fmt.Printf("Warning! Could not play music: %v\n", err)
						}
					} else { //If music is being played
						//If the music is paused
						if mix.PausedMusic() {
							//Resume the music
							mix.ResumeMusic()
						} else { //If the music is playing
							//Pause the music
							mix.PauseMusic()
						}
					}
					break
				case sdl.K_0:
					//Stop the music
					mix.HaltMusic()
					break
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

		err = gPromptTexture.Render(0, 0, nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("could not render prompt texture: %v", err)
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
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
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

	//Initialize SDL_mixer
	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 1024); err != nil {
		return fmt.Errorf("SDL_mixer could not initialize! SDL_mixer Error: %v", err)
	}

	return nil
}

func loadMedia() error {
	//Local error declaration
	var err error

	//Load prompt texture
	err = gPromptTexture.LoadFromFile("prompt.png")
	if err != nil {
		return fmt.Errorf("failed to load prompt texture: %v", err)
	}

	//Load music
	gMusic, err = mix.LoadMUS("beat.wav")
	if err != nil {
		return fmt.Errorf("failed to load beat music! SDL_mixer Error: %v", err)
	}

	//Load sound effects
	gScratch, err = mix.LoadWAV("scratch.wav")
	if err != nil {
		return fmt.Errorf("failed to load scratch sound effects! SDL_mixer Error: %v", err)
	}

	gHigh, err = mix.LoadWAV("high.wav")
	if err != nil {
		return fmt.Errorf("failed to load high sound effects! SDL_mixer Error: %v", err)
	}

	gMedium, err = mix.LoadWAV("medium.wav")
	if err != nil {
		return fmt.Errorf("failed to load medium sound effects! SDL_mixer Error: %v", err)
	}

	gLow, err = mix.LoadWAV("low.wav")
	if err != nil {
		return fmt.Errorf("failed to load low sound effects! SDL_mixer Error: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gPromptTexture.Free(); err != nil {
		return fmt.Errorf("could not free prompt texture: %v", err)
	}

	//Free the sound effects
	gScratch.Free()
	gHigh.Free()
	gMedium.Free()
	gLow.Free()
	gScratch = nil
	gHigh = nil
	gMedium = nil
	gLow = nil

	//Free the music
	gMusic.Free()
	gMusic = nil

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
	mix.Quit()
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
