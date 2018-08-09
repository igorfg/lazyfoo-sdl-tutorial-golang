package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	//Screen dimension constants
	screenWitdh  = 640
	screenHeight = 480

	//The dimensions of the level
	levelWidth  = 1280
	levelHeight = 960

	//Tile constants
	tileWidth        = 80
	tileHeight       = 80
	totalTiles       = 192
	totalTileSprites = 12

	//The different tile sprites
	tileRed         = 0
	tileGreen       = 1
	tileBlue        = 2
	tileCenter      = 3
	tileTop         = 4
	tileTopRight    = 5
	tileRight       = 6
	tileBottomRight = 7
	tileBottom      = 8
	tileBottomLeft  = 9
	tileLeft        = 10
	tileTopLeft     = 11
)

var (
	//The window we'll be rendering to
	gWindow *sdl.Window

	//The window renderer
	gRenderer *sdl.Renderer

	//Globally used font
	gFont *ttf.Font

	//Scene textures
	gDotTexture  LTexture
	gTileTexture LTexture
	gTileClips   [totalTileSprites]sdl.Rect
)

func main() {
	//Start up SDL and create window
	if err := initSDl(); err != nil {
		log.Fatalf("Could not init SDL: %v\n", err)
	}

	//The level tiles
	tileSet := make([]*Tile, totalTiles)

	//Load media
	if err := loadMedia(tileSet); err != nil {
		log.Fatalf("Could not load media: %v\n", err)
	}

	//Main loop flag
	var quit bool

	//Event handler
	var e sdl.Event

	//The dot that will be moving around on the screen
	dot := NewDot()

	//Level camera
	camera := &sdl.Rect{X: 0, Y: 0, W: screenWitdh, H: screenHeight}

	//While application is running
	for !quit {
		//Handle events on queue
		for e = sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			//User requests quit
			if e.GetType() == sdl.QUIT {
				quit = true
			}

			//Handle input for the dot
			dot.HandleEvent(e)
		}

		//Move the dot
		dot.Move(tileSet)
		dot.SetCamera(camera)

		//Clear screen
		err := gRenderer.SetDrawColor(255, 255, 255, 255)
		if err != nil {
			log.Fatalf("could not set draw color for renderer: %v", err)
		}
		err = gRenderer.Clear()
		if err != nil {
			log.Fatalf("could not clear renderer: %v", err)
		}

		//Render level
		for i := 0; i < totalTiles; i++ {
			if err = tileSet[i].Render(camera); err != nil {
				log.Fatal(err)
			}
		}

		//Render dot
		if err = dot.Render(camera); err != nil {
			log.Fatalf("%v\n", err)
		}

		//Update screen
		gRenderer.Present()
	}

	//Free resources and close SDL
	if err := close(tileSet); err != nil {
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

func loadMedia(tiles []*Tile) error {
	var err error

	//Load dot texture
	if err = gDotTexture.LoadFromFile("dot.bmp"); err != nil {
		return fmt.Errorf("Failed to load dot texture: %v", err)
	}

	//Load tile texture
	if err = gTileTexture.LoadFromFile("tiles.png"); err != nil {
		return fmt.Errorf("Failed to load tile set texture: %v", err)
	}

	//Load tile map
	if err = setTiles(tiles); err != nil {
		return fmt.Errorf("Failed to load tile set: %v", err)
	}

	return nil
}

func close(tiles []*Tile) error {
	//Free loaded images
	if err := gDotTexture.Free(); err != nil {
		return fmt.Errorf("could not free dot texture: %v", err)
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

//Box collision detector
func checkCollision(a, b sdl.Rect) bool {
	//The sides of the rectangles
	var (
		leftA, leftB     int32
		rightA, rightB   int32
		topA, topB       int32
		bottomA, bottomB int32
	)

	//Calculate the sides of rect A
	leftA = a.X
	rightA = a.X + a.W
	topA = a.Y
	bottomA = a.Y + a.H

	//Calculate the sides of rect B
	leftB = b.X
	rightB = b.X + b.W
	topB = b.Y
	bottomB = b.Y + b.H

	//If any of the sides from A are outside of B
	if bottomA <= topB {
		return false
	}
	if topA >= bottomB {
		return false
	}
	if rightA <= leftB {
		return false
	}
	if leftA >= rightB {
		return false
	}

	//If none of the sides from A are outside of B
	return true
}

//Checks collision box against set of tiles
func touchesWall(box sdl.Rect, tiles []*Tile) bool {
	//Go through tiles
	for i := 0; i < totalTiles; i++ {
		//If the tile is a wall type tile
		if tiles[i].MType() >= tileCenter && tiles[i].MType() <= tileTopLeft {
			//If collision box touches the wall tile
			if checkCollision(box, tiles[i].MBox()) {
				return true
			}
		}
	}

	//If no wall tiles were touched
	return false
}

//Sets tiles from tile map
func setTiles(tiles []*Tile) error {
	//The tile offsets
	var x, y int32

	//Open the map
	mapFile, err := os.Open("lazy.map")
	if err != nil {
		return fmt.Errorf("Unable to load map file: %v", err)
	}
	defer mapFile.Close()

	//Scanner that will be used to read the tile numbers
	scanner := bufio.NewScanner(mapFile)
	scanner.Split(bufio.ScanWords)

	//Initialize the tiles
	for i := 0; i < totalTiles; i++ {
		//Determines what kind of tile will be made
		tileType := -1

		//Read tile from map file
		scanner.Scan()
		tileType, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("Error loading map: Unexpected EOF")
		}

		//If the number is valid tile number
		if tileType >= 0 && tileType < totalTileSprites {
			tiles[i] = NewTile(x, y, tileType)
		} else {
			return fmt.Errorf("Error loading map: Invalid tile type at %d", i)
		}

		//Move to next tile spot
		x += tileWidth

		//If we've gone too far
		if x >= levelWidth {
			//Move back
			x = 0

			//Move to the next row
			y += tileHeight
		}
	}

	//Clip the sprite sheet
	gTileClips[tileRed].X = 0
	gTileClips[tileRed].Y = 0
	gTileClips[tileRed].W = tileWidth
	gTileClips[tileRed].H = tileHeight

	gTileClips[tileGreen].X = 0
	gTileClips[tileGreen].Y = 80
	gTileClips[tileGreen].W = tileWidth
	gTileClips[tileGreen].H = tileHeight

	gTileClips[tileBlue].X = 0
	gTileClips[tileBlue].Y = 160
	gTileClips[tileBlue].W = tileWidth
	gTileClips[tileBlue].H = tileHeight

	gTileClips[tileTopLeft].X = 80
	gTileClips[tileTopLeft].Y = 0
	gTileClips[tileTopLeft].W = tileWidth
	gTileClips[tileTopLeft].H = tileHeight

	gTileClips[tileLeft].X = 80
	gTileClips[tileLeft].Y = 80
	gTileClips[tileLeft].W = tileWidth
	gTileClips[tileLeft].H = tileHeight

	gTileClips[tileBottomLeft].X = 80
	gTileClips[tileBottomLeft].Y = 160
	gTileClips[tileBottomLeft].W = tileWidth
	gTileClips[tileBottomLeft].H = tileHeight

	gTileClips[tileTop].X = 160
	gTileClips[tileTop].Y = 0
	gTileClips[tileTop].W = tileWidth
	gTileClips[tileTop].H = tileHeight

	gTileClips[tileCenter].X = 160
	gTileClips[tileCenter].Y = 80
	gTileClips[tileCenter].W = tileWidth
	gTileClips[tileCenter].H = tileHeight

	gTileClips[tileBottom].X = 160
	gTileClips[tileBottom].Y = 160
	gTileClips[tileBottom].W = tileWidth
	gTileClips[tileBottom].H = tileHeight

	gTileClips[tileTopRight].X = 240
	gTileClips[tileTopRight].Y = 0
	gTileClips[tileTopRight].W = tileWidth
	gTileClips[tileTopRight].H = tileHeight

	gTileClips[tileRight].X = 240
	gTileClips[tileRight].Y = 80
	gTileClips[tileRight].W = tileWidth
	gTileClips[tileRight].H = tileHeight

	gTileClips[tileBottomRight].X = 240
	gTileClips[tileBottomRight].Y = 160
	gTileClips[tileBottomRight].W = tileWidth
	gTileClips[tileBottomRight].H = tileHeight

	return nil
}
