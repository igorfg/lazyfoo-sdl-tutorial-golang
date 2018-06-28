//The mouse button
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//LButtonSprite defines enum type for LButton
type LButtonSprite int

//Button constants
const (
	buttonWidth  = 300
	buttonHeight = 200
	totalButtons = 4
)

//Button sprite enum
const (
	buttonSpriteMouseOut LButtonSprite = iota
	buttonSpriteMouseOverMotion
	buttonSpriteMouseDown
	buttonSpriteMouseUp
	buttonSpriteTotal
)

//LButton is the struct that loads a button
type LButton struct {
	//Top left position
	mPosition sdl.Point

	//Current used global sprite
	mCurrentSprite LButtonSprite
}

//NewLButton initializes internal variables
func (lb *LButton) NewLButton() {
	lb.mPosition.X = 0
	lb.mPosition.Y = 0

	lb.mCurrentSprite = buttonSpriteMouseOut
}

//SetPosition sets top left position
func (lb *LButton) SetPosition(x, y int32) {
	lb.mPosition.X = x
	lb.mPosition.Y = y
}

//HandleEvent handles mouse event
func (lb *LButton) HandleEvent(e sdl.Event) {
	//If mouse event happened
	if e.GetType() == sdl.MOUSEMOTION || e.GetType() == sdl.MOUSEBUTTONDOWN || e.GetType() == sdl.MOUSEBUTTONUP {
		//Get mouse position
		x, y, _ := sdl.GetMouseState()

		//Check if mouse is in button
		inside := true

		//Mouse is left of the button
		if x < lb.mPosition.X {
			inside = false
		} else if x > lb.mPosition.X+buttonWidth { //Mouse is right of the button
			inside = false
		} else if y < lb.mPosition.Y { //Mouse above the button
			inside = false
		} else if y > lb.mPosition.Y+buttonHeight { //Mouse below the button
			inside = false
		}

		//Mouse is outside button
		if !inside {
			lb.mCurrentSprite = buttonSpriteMouseOut
		} else { //Mouse is inside button
			//Set mouse over sprite
			switch e.GetType() {
			case sdl.MOUSEMOTION:
				lb.mCurrentSprite = buttonSpriteMouseOverMotion
				break
			case sdl.MOUSEBUTTONDOWN:
				lb.mCurrentSprite = buttonSpriteMouseDown
				break
			case sdl.MOUSEBUTTONUP:
				lb.mCurrentSprite = buttonSpriteMouseUp
				break
			}
		}
	}
}

//Render shows button sprite
func (lb *LButton) Render() error {
	//Show current button sprite
	err := gButtonSpriteSheetTexture.Render(lb.mPosition.X, lb.mPosition.Y, &gSpriteClips[lb.mCurrentSprite], 0, nil, sdl.FLIP_NONE)
	if err != nil {
		return fmt.Errorf("could not render button sprite: %v", err)
	}
	return nil
}
