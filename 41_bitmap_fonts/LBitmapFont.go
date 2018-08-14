package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

//LBitmapFont is our bitmap font
type LBitmapFont struct {
	//The font texture
	mBitmap *LTexture

	//The individual characters in the surface
	mChars [256]sdl.Rect

	//Spacing variables
	mNewLine int32
	mSpace   int32
}

//BuildFont generates the font
func (bmf *LBitmapFont) BuildFont(bitmap *LTexture) error {
	//Lock pixels for access
	if err := bitmap.LockTexture(); err != nil {
		return fmt.Errorf("unable to lock bitmap font texture: %v", err)
	}

	//Set the background color
	bgColor := bitmap.GetPixel32(0, 0)

	//Set the cell dimensions
	cellW := bitmap.MWidth() / 16
	cellH := bitmap.MHeight() / 16

	//New line variables
	top := cellH
	baseA := cellH

	//The current character we are setting
	currentChar := 0

	//Go through the cell rows
	for rows := int32(0); rows < 16; rows++ {
		//Go through the cell columns
		for cols := int32(0); cols < 16; cols++ {
			//Set the character offset
			bmf.mChars[currentChar].X = cellW * cols
			bmf.mChars[currentChar].Y = cellH * rows

			//Set the dimension of the caracter
			bmf.mChars[currentChar].W = cellW
			bmf.mChars[currentChar].H = cellH

			//Find the left side
			//Go through pixel columns
			for pCol := int32(0); pCol < cellW; pCol++ {
				//Go through pixel rows
				for pRow := int32(0); pRow < cellH; pRow++ {
					//Get the pixel offsets
					pX := (cellW * cols) + pCol
					pY := (cellH * rows) + pRow

					//If a non colorkey is found
					if bitmap.GetPixel32(int(pX), int(pY)) != bgColor {
						//Set the x offset
						bmf.mChars[currentChar].X = pX

						//Break the loops
						pCol = cellW
						pRow = cellH
					}
				}
			}

			//Find the right side
			//Go through pixel columns
			for pColW := cellW - 1; pColW >= 0; pColW-- {
				//Go through pixel rows
				for pRowW := int32(0); pRowW < cellH; pRowW++ {
					//Get the pixel offsets
					pX := (cellW * cols) + pColW
					pY := (cellH * rows) + pRowW

					//If a non colorkey pixel is found
					if bitmap.GetPixel32(int(pX), int(pY)) != bgColor {
						//Set the width
						bmf.mChars[currentChar].W = (pX - bmf.mChars[currentChar].X) + 1

						//Break the loops
						pColW = -1
						pRowW = cellH
					}
				}
			}

			//Find top
			//Go through pixel rows
			for pRow := int32(0); pRow < cellH; pRow++ {
				//Go through pixel columns
				for pCol := int32(0); pCol < cellW; pCol++ {
					//Get the pixel offsets
					pX := (cellW * cols) + pCol
					pY := (cellH * rows) + pRow

					//If a non colorkey pixel is found
					if bitmap.GetPixel32(int(pX), int(pY)) != bgColor {
						//If a new top is found
						if pRow < top {
							top = pRow
						}

						//Break the loops
						pCol = cellW
						pRow = cellH
					}
				}
			}

			//Find bottom of A
			if currentChar == 'A' {
				//Go through pixel rows
				for pRow := cellH - 1; pRow >= 0; pRow-- {
					//Go through pixel rows
					for pCol := int32(0); pCol < cellW; pCol++ {
						//Get the pixel offsets
						pX := (cellW * cols) + pCol
						pY := (cellH * rows) + pRow

						//If a non colorkey pixel is found
						if bitmap.GetPixel32(int(pX), int(pY)) != bgColor {
							//Bottom of A is found
							baseA = pRow

							//Break the loops
							pCol = cellW
							pRow = -1
						}
					}
				}
			}
			//Go to the next character
			currentChar++
		}
	}

	//Calculate space
	bmf.mSpace = cellW / 2

	//Calculate new line
	bmf.mNewLine = baseA - top

	//Lop off excess top pixels
	for i := 0; i < 256; i++ {
		bmf.mChars[i].Y += top
		bmf.mChars[i].H -= top
	}

	if err := bitmap.UnlockTexture(); err != nil {
		return fmt.Errorf("could not unlock bitmap font texture: %v", err)
	}
	bmf.mBitmap = bitmap

	return nil
}

//RenderText renders the text
func (bmf *LBitmapFont) RenderText(x, y int32, text string) error {
	//If the font has been built
	if bmf.mBitmap != nil {
		//Temp offsets
		curX := x
		curY := y

		//Go through the text
		for i := 0; i < len(text); i++ {
			//If the current character is a space
			if text[i] == ' ' {
				//Move over
				curX += bmf.mSpace
			} else if text[i] == '\n' { //If the current character is a newline
				//Move down
				curY += bmf.mNewLine

				//Move back
				curX = x
			} else {
				//Get the ASCII value of the character
				ascii := text[i]

				//Show the character
				err := bmf.mBitmap.Render(curX, curY, &bmf.mChars[ascii], 0, nil, sdl.FLIP_NONE)
				if err != nil {
					return fmt.Errorf("could not render bitmap font character: %v", err)
				}

				//Move over the width of the character with one pixel padding
				curX += bmf.mChars[ascii].W + 1
			}
		}
	}

	return nil
}
