package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"

	"github.com/veandco/go-sdl2/sdl"
)

//DataStream is a test animation stream
type DataStream struct {
	//Internal data
	mImages       [4]*sdl.Surface
	mCurrentImage int
	mDelayFrames  int
}

//NewDataStream initializes internals
func NewDataStream() *DataStream {
	return &DataStream{mDelayFrames: 4}
}

//LoadMedia loads initial data
func (ds *DataStream) LoadMedia() error {
	for i := 0; i < 4; i++ {
		path := ""
		path = fmt.Sprintf("foo_walk_%d.png", i)

		loadedSurface, err := img.Load(path)
		if err != nil {
			return fmt.Errorf("unable to load %s: %v", path, err)
		}

		if ds.mImages[i], err = loadedSurface.ConvertFormat(sdl.PIXELFORMAT_RGBA8888, 0); err != nil {
			return fmt.Errorf("could not convert surface format: %v", err)
		}

		loadedSurface.Free()
	}

	return nil
}

//Free is the deallocator
func (ds *DataStream) Free() {
	for i := 0; i < 4; i++ {
		ds.mImages[i].Free()
	}
}

//GetBuffer gets current frame data
func (ds *DataStream) GetBuffer() []byte {
	ds.mDelayFrames--
	if ds.mDelayFrames == 0 {
		ds.mCurrentImage++
		ds.mDelayFrames = 4
	}

	if ds.mCurrentImage == 4 {
		ds.mCurrentImage = 0
	}

	return ds.mImages[ds.mCurrentImage].Pixels()
}
