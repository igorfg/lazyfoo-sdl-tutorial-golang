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
	gDotTexture        LTexture
	gPromptTextTexture LTexture
	gInputTextTexture  LTexture
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

	//Set color text as black
	var textColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//The current input text
	var inputText = "Some Text"
	err := gInputTextTexture.loadFromRenderedText(inputText, textColor)
	if err != nil {
		log.Fatalf("could not load input text texture: %v\n", err)
	}

	//Enable text input
	sdl.StartTextInput()

	//While application is running
	for !quit {
		//The rendered text flag
		renderText := false

		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			} else if e.GetType() == sdl.KEYDOWN { //Special key input
				//Getting the keyboard event from the Event interface
				k := e.(*sdl.KeyboardEvent)

				//Handle backspace
				if k.Keysym.Sym == sdl.K_BACKSPACE && len(inputText) > 0 {
					//Lop off character
					inputText = inputText[:len(inputText)-1]
					renderText = true
				} else if k.Keysym.Sym == sdl.K_c && sdl.GetModState()&sdl.KMOD_CTRL != 0 { //Handle copy
					if err = sdl.SetClipboardText(inputText); err != nil {
						log.Fatalf("could not set clipboard text: %v\n", err)
					}
				} else if k.Keysym.Sym == sdl.K_v && sdl.GetModState()&sdl.KMOD_CTRL != 0 { //Handle paste
					if inputText, err = sdl.GetClipboardText(); err != nil {
						log.Fatalf("could not get clipboard text: %v\n", err)
					}
					renderText = true
				}
			} else if e.GetType() == sdl.TEXTINPUT { //Special text input event
				//Getting the text input event from the Event interface
				t := e.(*sdl.TextInputEvent)

				//Not copy or pasting
				if !((t.Text[0] == 'c' || t.Text[0] == 'C') &&
					(t.Text[0] == 'v' || t.Text[0] == 'V') &&
					(sdl.GetModState() == sdl.KMOD_CTRL)) {

					//Append character
					inputText = string(append([]byte(inputText), t.Text[0]))
					renderText = true
				}
			}
		}

		//Rerender text if needed
		if renderText {
			//Text is not empty
			if inputText != "" {
				//Render new text
				err := gInputTextTexture.loadFromRenderedText(inputText, textColor)
				if err != nil {
					log.Fatalf("could not render input texture: %v\n", err)
				}
			} else { //Text is empty
				//Render space texture becase SDL_TTF does not render empty strings
				err := gInputTextTexture.loadFromRenderedText(" ", textColor)
				if err != nil {
					log.Fatalf("could not render input texture: %v\n", err)
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
		err = gInputTextTexture.Render((screenWitdh-gInputTextTexture.GetWidth())/2,
			gPromptTextTexture.GetHeight(), nil, 0, nil, sdl.FLIP_NONE)
		if err != nil {
			log.Fatalf("coud not render input texture: %v\n", err)
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
	//Render the prompt
	textColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//Load prompt font texture
	err = gPromptTextTexture.loadFromRenderedText("Enter text: ", textColor)
	if err != nil {
		return fmt.Errorf("Failed to load prompt text: %v", err)
	}

	return nil
}

func close() error {
	//Free loaded images
	if err := gPromptTextTexture.Free(); err != nil {
		return fmt.Errorf("could not free prompt texture: %v", err)
	}
	if err := gInputTextTexture.Free(); err != nil {
		return fmt.Errorf("could not free input texture: %v", err)
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
