package display

import "image"

// ScreenGrabber TODO
type Grabber interface {
	Start()
	Frames() <-chan *image.RGBA
	Stop()
	Fps() int
	Bitrate() int
	Screen() *Screen
	SendEvent([]byte)
}

// Event TODO
type Event interface {
	getEventType() int
}

// Screen TODO
type Screen struct {
	Index  int
	Bitrate int 
	FrameRate int
	Name   string
	Bounds image.Rectangle
}

// Service TODO
type Service interface {
	CreateGrabber(screen Screen, fps int, bitrate int) (Grabber, error)
	Screens() ([]Screen, error)
}

