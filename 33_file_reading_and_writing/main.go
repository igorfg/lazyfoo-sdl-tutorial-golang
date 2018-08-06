package main

import (
	"fmt"
	"log"
	"strconv"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	//Screen dimension constants
	screenWitdh  = 640
	screenHeight = 480

	//Number of data integers
	totalData = 10
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gPromptTextTexture LTexture
	gDataTextures      [totalData]LTexture

	//Data points
	gData [totalData]int32
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

	//Rendering colors
	var textColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	var highlightColor = sdl.Color{R: 255, G: 0, B: 0, A: 255}

	//Current input point
	var currentData = 0

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.KEYDOWN {
				switch (e.(*sdl.KeyboardEvent)).Keysym.Sym {
				//Previous data entry
				case sdl.K_UP:
					//Rerender previous entry input point
					err := gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), textColor)
					if err != nil {
						log.Fatalf("could not render previous input point: %v", err)
					}

					currentData--
					if currentData < 0 {
						currentData = totalData - 1
					}

					//Rerender current entry input point
					err = gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), highlightColor)
					if err != nil {
						log.Fatalf("could not render current input point: %v", err)
					}
					break
				case sdl.K_DOWN:
					//Rerender previous entry input point
					err := gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), textColor)
					if err != nil {
						log.Fatalf("could not render previous input point: %v", err)
					}

					currentData++
					if currentData == totalData {
						currentData = 0
					}

					//Rerender current entry input point
					err = gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), highlightColor)
					if err != nil {
						log.Fatalf("could not render current input point: %v", err)
					}
					break
				case sdl.K_LEFT:
					gData[currentData]--
					err := gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), highlightColor)
					if err != nil {
						log.Fatalf("could not render after decrementing current data: %v", err)
					}
					break
				case sdl.K_RIGHT:
					gData[currentData]++
					err := gDataTextures[currentData].loadFromRenderedText(strconv.Itoa(int(gData[currentData])), highlightColor)
					if err != nil {
						log.Fatalf("could not render after incrementing current data: %v", err)
					}
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

		//Render text textures
		err = gPromptTextTexture.Render((screenWitdh-gPromptTextTexture.GetWidth())/2,
			0, nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("coud not render prompt texture: %v\n", err)
		}

		for i := 0; i < totalData; i++ {
			err = gDataTextures[i].Render((screenWitdh-gDataTextures[i].GetWidth())/2,
				gPromptTextTexture.GetHeight()+gDataTextures[0].GetHeight()*int32(i),
				nil, 0, nil, sdl.FLIP_NONE)
			if err != nil {
				log.Fatalf("could not render data texture %d: %v\n", i, err)
			}
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

	//Text rendering colors
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	highlightColor := sdl.Color{R: 255, G: 0, B: 0, A: 255}

	//Open the font
	gFont, err = ttf.OpenFont("lazy.ttf", 28)
	if err != nil {
		return fmt.Errorf("Failed to load lazy font! SDL_ttf Error: %v", err)
	}

	//Render the prompt
	err = gPromptTextTexture.loadFromRenderedText("Enter data:", textColor)
	if err != nil {
		return fmt.Errorf("Failed to render prompt text: %v", err)
	}

	//Open the file for reading in binary
	file := sdl.RWFromFile("nums.bin", "r+b")

	//File does not exist
	if file == nil {
		fmt.Printf("Warning: unable to open file! SDL_error: %s\n", sdl.GetError())
		file = sdl.RWFromFile("nums.bin", "w+b")

		if file != nil {
			fmt.Println("New file created!")

			//Initialize data
			for i := 0; i < totalData; i++ {
				file.RWwrite(unsafe.Pointer(&gData[i]), 2, 1)
			}

			//Close file handler
			err = file.RWclose()
			if err != nil {
				fmt.Println("Warning: could not close file: ", err)
			}
		} else {
			return fmt.Errorf("Error: Unable to create file! SDL_Error: %s", sdl.GetError())
		}
	} else { //File exists
		//Load data
		fmt.Println("Reading file...!")
		for i := 0; i < totalData; i++ {
			file.RWread(unsafe.Pointer(&gData[i]), uint(unsafe.Sizeof(gData[i])), 1)
		}

		//Close file handler
		err = file.RWclose()
		if err != nil {
			fmt.Println("Warning: could not close file: ", err)
		}
	}

	//Initialize data textures
	err = gDataTextures[0].loadFromRenderedText(strconv.Itoa(int(gData[0])), highlightColor)
	if err != nil {
		return fmt.Errorf("could not load %d gData texture on initialization: %v", 0, err)
	}
	for i := 1; i < totalData; i++ {
		err = gDataTextures[i].loadFromRenderedText(strconv.Itoa(int(gData[i])), textColor)
		if err != nil {
			return fmt.Errorf("could not load %d gData texture: %v", i, err)
		}
	}

	return nil
}

func close() error {
	//Open data for writing
	file := sdl.RWFromFile("nums.bin", "w+b")

	if file != nil {
		//Save data
		for i := 0; i < totalData; i++ {
			file.RWwrite(unsafe.Pointer(&gData[i]), 2, 1)
		}
		//Close file handler
		err := file.RWclose()
		if err != nil {
			fmt.Println("Warning: could not close file: ", err)
		}
	} else {
		return fmt.Errorf("Error: Unable to save file! SDL_Error: %s", sdl.GetError())
	}

	//Free loaded images
	if err := gPromptTextTexture.Free(); err != nil {
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
