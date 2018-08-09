package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

//LTimer is the application time based timer
type LTimer struct {
	//The clock time when the timer started
	mStartTicks uint32

	//The tics stored when the timer was paused
	mPausedTicks uint32

	//The time status
	mPaused  bool
	mStarted bool
}

//Start starts the clock
func (lt *LTimer) Start() {
	//Start the timer
	lt.mStarted = true

	//Unpause the timer
	lt.mPaused = false

	//Get the current clock time
	lt.mStartTicks = sdl.GetTicks()
	lt.mPausedTicks = 0
}

//Stop stops the clock
func (lt *LTimer) Stop() {
	//Stop te timer
	lt.mStarted = false

	//Unpause the timer
	lt.mPaused = false

	//Clear tick variables
	lt.mStartTicks = 0
	lt.mPausedTicks = 0
}

//Pause pauses the clock
func (lt *LTimer) Pause() {
	//If the timer is running and isn't already paused
	if lt.mStarted && !lt.mPaused {
		//Pause the timer
		lt.mPaused = true

		//Calculate the paused ticks
		lt.mPausedTicks = sdl.GetTicks() - lt.mStartTicks
		lt.mStartTicks = 0
	}
}

//Unpause unpauses the clock
func (lt *LTimer) Unpause() {
	//If the timer is running and paused
	if lt.mStarted && lt.mPaused {
		//Unpause the timer
		lt.mPaused = false

		//Reset the starting ticks
		lt.mStartTicks = sdl.GetTicks() - lt.mPausedTicks

		//Reset the paused ticks
		lt.mPausedTicks = 0
	}
}

//GetTicks gets the timer's time
func (lt *LTimer) GetTicks() uint32 {
	//The actual timer line
	var time uint32

	//If the timer is running
	if lt.mStarted {
		//If the timer is paused
		if lt.mPaused {
			//Return the number of ticks when the timer was paused
			time = lt.mPausedTicks
		} else {
			//Return the current time minus the start time
			time = sdl.GetTicks() - lt.mStartTicks
		}
	}

	return time
}

//IsStarted checks if the the timer started
func (lt *LTimer) IsStarted() bool {
	//Timer is running and paused or unpaused
	return lt.mStarted
}

//IsPaused checks if the timer was paused
func (lt *LTimer) IsPaused() bool {
	//Timer is running and paused
	return lt.mPaused && lt.mStarted
}
