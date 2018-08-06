package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

//LTexture is a sdl.Texture wrapper class
type LTexture struct {
	//The actual hardware texture
	mTexture *sdl.Texture

	//image dimensions
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
	err := lt.Free()
	if err != nil {
		return fmt.Errorf("could not free LTexture: %v", err)
	}

	//The final texture
	var newTexture *sdl.Texture

	//Load image at specified path
	loadedSurface, err := img.Load(path)
	if err != nil {
		return fmt.Errorf("could not load image %v! SDL_image Error: %v", path, err)
	}

	//Color key image
	err = loadedSurface.SetColorKey(true, sdl.MapRGB(loadedSurface.Format, 0, 255, 255))
	if err != nil {
		return fmt.Errorf("could not set color key: %v", err)
	}

	//Create texture from surface pixels
	newTexture, err = gRenderer.CreateTextureFromSurface(loadedSurface)
	if err != nil {
		return fmt.Errorf("could not create texture from %v pixels: %v", path, err)
	}

	//Get image dimensions
	lt.mWidth = loadedSurface.W
	lt.mHeight = loadedSurface.H

	//Get rid of old loaded surface
	loadedSurface.Free()

	//Assing new texture
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

//GetWidth gets image width
func (lt *LTexture) GetWidth() int32 {
	return lt.mWidth
}

//GetHeight gets image height
func (lt *LTexture) GetHeight() int32 {
	return lt.mHeight
}
