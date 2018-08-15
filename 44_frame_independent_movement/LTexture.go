package main

import (
	"encoding/binary"
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// TextureAccess is an enumartion of texture access patters
// (https://wiki.libsdl.org/SDL_TextureAccess)
type TextureAccess int

//LTexture is a sdl.Texture wrapper class
type LTexture struct {
	//The actual hardware texture
	mTexture *sdl.Texture
	mPixels  []byte
	mPitch   int

	//Image dimensions
	mWidth  int32
	mHeight int32
}

//NewLTexture initializes variables
func NewLTexture() *LTexture {
	//Initialize
	return &LTexture{mTexture: nil, mWidth: 0, mHeight: 0}
}

//LoadFromFile loads image at specified path
func (lt *LTexture) LoadFromFile(path string) error {
	//Get rid of preexisting texture
	if err := lt.Free(); err != nil {
		return fmt.Errorf("could not free texture: %v", err)
	}

	//The final texture
	var newTexture *sdl.Texture

	//Load image at specified path
	loadedSurface, err := img.Load(path)
	if err != nil {
		return fmt.Errorf("could not load image %v! SDL_image Error: %v", path, err)
	}

	//Convert surface to display format
	formattedSurface, err := loadedSurface.ConvertFormat(sdl.PIXELFORMAT_ABGR8888, 0)
	if err != nil {
		return fmt.Errorf("could not convert surface to display format: %v", err)
	}

	//Create blank streamable texture
	newTexture, err = gRenderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING,
		formattedSurface.W, formattedSurface.H)
	if err != nil {
		return fmt.Errorf("could not create blank texture: %v", err)
	}

	//Enable blending on texture
	if err := newTexture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return fmt.Errorf("could not set texture's blend mode: %v", err)
	}

	//Lock texture for manipulation
	lt.mPixels, lt.mPitch, err = newTexture.Lock(&formattedSurface.ClipRect)
	if err != nil {
		return fmt.Errorf("could not lock texture: %v", err)
	}

	//Copy loaded/formatted surface pixels
	copy(lt.mPixels, formattedSurface.Pixels())

	//Get image dimensions
	lt.mWidth = formattedSurface.W
	lt.mHeight = formattedSurface.H

	//Get pixel data in editable format
	bytePixels := lt.mPixels
	pixels := make([]uint32, len(bytePixels)/4)
	for i := range pixels {
		//Assuming little endian
		pixels[i] = uint32(binary.LittleEndian.Uint32(bytePixels[i*4 : (i+1)*4]))
	}

	pixelCount := (lt.mPitch / 4) * int(lt.mHeight)

	//Map colors
	colorKey := sdl.MapRGB(formattedSurface.Format, 0, 255, 255)
	transparent := sdl.MapRGBA(formattedSurface.Format, 0, 255, 0, 0)

	//Color key pixels
	for i := 0; i < pixelCount; i++ {
		if pixels[i] == colorKey {
			binary.LittleEndian.PutUint32(bytePixels[i*4:(i*4)+4], transparent)
		}
	}

	//Unlock texture to update
	newTexture.Unlock()
	lt.mPixels = nil

	//Get rid of old formatted surface
	formattedSurface.Free()

	//Get rid of old loaded surface
	loadedSurface.Free()

	//Return no errors
	lt.mTexture = newTexture
	return nil
}

//LoadFromRenderedText creates image from font string
func (lt *LTexture) loadFromRenderedText(textureText string, textColor sdl.Color) error {
	var textSurface *sdl.Surface

	//Get rid of preexisting texture
	lt.Free()

	//Render text surface

	textSurface, err := gFont.RenderUTF8Solid(textureText, textColor)
	if err != nil {
		return fmt.Errorf("unable to render text surface! SDL_ttf Error: %v", err)
	}

	//Create texture from surface pixels
	lt.mTexture, err = gRenderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		return fmt.Errorf("unable to create texture from rendered text! SDL Error: %v", err)
	}

	//Get image dimensions
	lt.mWidth = textSurface.W
	lt.mHeight = textSurface.H

	//Get rid of old surface
	textSurface.Free()

	return nil
}

//Free deallocates memory
func (lt *LTexture) Free() error {
	//Free texture if it exists
	if lt.mTexture != nil {
		err := lt.mTexture.Destroy()
		if err != nil {
			return fmt.Errorf("could not destroy LTexture: %v", err)
		}

		lt.mTexture = nil
		lt.mWidth = 0
		lt.mHeight = 0
	}
	return nil
}

//SetColor sets color modulation
func (lt *LTexture) SetColor(red, green, blue uint8) error {
	//Modulate texture
	err := lt.mTexture.SetColorMod(red, green, blue)
	if err != nil {
		return fmt.Errorf("could not set color mod for texture: %v", err)
	}
	return nil
}

//SetBlendMode sets blending
func (lt *LTexture) SetBlendMode(blending sdl.BlendMode) error {
	//Set blending function
	err := lt.mTexture.SetBlendMode(blending)
	if err != nil {
		return fmt.Errorf("could not set blend mode: %v", err)
	}

	return nil
}

//SetAlpha sets alpha modulation
func (lt *LTexture) SetAlpha(alpha uint8) error {
	//Modulate texture alpha
	err := lt.mTexture.SetAlphaMod(alpha)
	if err != nil {
		return fmt.Errorf("could not set alpha mod: %v", err)
	}

	return nil
}

//Render renders texture at given point
func (lt *LTexture) Render(x, y int32, clip *sdl.Rect, angle float64, center *sdl.Point, flip sdl.RendererFlip) error {
	//Set rendering space and render to screen
	renderQuad := sdl.Rect{X: x, Y: y, W: lt.mWidth, H: lt.mHeight}

	//Set clip rendering dimenisions
	if clip != nil {
		renderQuad.W = clip.W
		renderQuad.H = clip.H
	}

	//Render to screen
	err := gRenderer.CopyEx(lt.mTexture, clip, &renderQuad, angle, center, flip)
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

//LockTexture locks texture for pixel manipulation
func (lt *LTexture) LockTexture() error {
	var err error

	//Texture is already locked
	if lt.mPixels != nil {
		return fmt.Errorf("texture is already locked")
	}

	lt.mPixels, lt.mPitch, err = lt.mTexture.Lock(nil)
	if err != nil {
		return fmt.Errorf("unable to lock texture: %v", err)
	}

	return nil
}

//UnlockTexture unlocks texture for pixel manipulation
func (lt *LTexture) UnlockTexture() error {
	//Texture is not locked
	if lt.mPixels == nil {
		return fmt.Errorf("texture is not locked")
	}

	//Unlock texture
	lt.mTexture.Unlock()
	lt.mPixels = nil
	lt.mPitch = 0

	return nil
}

//CreateBlank creates blank texture
func (lt *LTexture) CreateBlank(width, height int32, access TextureAccess) error {
	var err error

	//Create unitialized texture
	lt.mTexture, err = gRenderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, int(access), width, height)
	if err != nil {
		return fmt.Errorf("unable to create blank texture: %v", err)
	}

	lt.mWidth = width
	lt.mHeight = height

	return nil
}

//CopyPixels copies pixels
func (lt *LTexture) CopyPixels(pixels []byte) {
	//Texture is locked
	if lt.mPixels != nil {
		//Copy to locked pixels
		copy(lt.mPixels, pixels)
	}
}

//SetAsRenderTarget sets self as render target
func (lt *LTexture) SetAsRenderTarget() error {
	//Make self render target
	if err := gRenderer.SetRenderTarget(lt.mTexture); err != nil {
		return fmt.Errorf("could not set render target: %v", err)
	}
	return nil
}

//MWidth gets image width
func (lt *LTexture) MWidth() int32 {
	return lt.mWidth
}

//MHeight gets image height
func (lt *LTexture) MHeight() int32 {
	return lt.mHeight
}

//MPixels gets texture pixels' start address
func (lt *LTexture) MPixels() []byte {
	return lt.mPixels
}

//MPitch gets texture's pitch
func (lt *LTexture) MPitch() int {
	return lt.mPitch
}

//GetPixel32 at exact (x,y) coordinate
func (lt *LTexture) GetPixel32(x, y int) uint32 {
	//Convert the pixel to 32 bit
	position := y*lt.mPitch + x*4

	return uint32(binary.LittleEndian.Uint32(lt.mPixels[position : position+4]))
}
