// +build linux

package display

import (
	"image"
	"time"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/wadahana/memu/log"
)

// LinuxProvider implements the rdisplay.Service interface for XServer
type LinuxProvider struct{}

// LinuxGrabber captures video from a X server
type LinuxGrabber struct {
	fps     int
	bitrate int
	screen Screen
	frames chan *image.RGBA
	stop   chan struct {}
	events chan Event
	
}

// CreateGrabber Creates an screen capturer for the X server
func (*LinuxProvider) CreateGrabber(screen Screen, fps int, bitrate int) (Grabber, error) {
	return &LinuxGrabber {
		screen:  screen,
		fps:     fps,
		bitrate: bitrate,
		frames:  make(chan *image.RGBA),
		stop:    make(chan struct{}),
		events:  make(chan Event),
	}, nil
}

// Screens Returns the available screens to capture
func (x *LinuxProvider) Screens() ([]Screen, error) {
	numScreens := screenshot.NumActiveDisplays()
	screens := make([]Screen, numScreens)
	for i := 0; i < numScreens; i++ {
		screens[i] = Screen{
			Index:  i,
			Name: fmt.Sprintf("Screen_%d", i),
			Bounds: screenshot.GetDisplayBounds(i),
		}
	}
	return screens, nil
}

// Frames returns a channel that will receive an image stream
func (g *LinuxGrabber) Frames() <-chan *image.RGBA {
	return g.frames
}

// Start initiates the screen capture loop
func (g *LinuxGrabber) Start() {
	delta := time.Duration(1000/g.fps) * time.Millisecond
	go func() {
		for {
			startedAt := time.Now()
			select {
			case <-g.stop:
				close(g.frames)
				close(g.events)
				return
			case e := <-g.events:
				if e.getEventType() == 1 {
					mouseEvent := e.(*MouseEvent);
					log.Infof("MouseEvent: %02x,%02x, %0.4f, %0.4f\n", mouseEvent.getEventType(), mouseEvent.mouseType, mouseEvent.x, mouseEvent.y);
				} else {
					log.Info("Unknown Event");
				}
			default:
				img, err := screenshot.CaptureRect(g.screen.Bounds)
				if err != nil {
					return
				}
				g.frames <- img
				ellapsed := time.Now().Sub(startedAt)
				sleepDuration := delta - ellapsed
				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}
		}
	}()
}

// Stop sends a stop signal to the capture loop
func (g *LinuxGrabber) Stop() {
	close(g.stop)
}

// Screen returns a pointer to the screen we're capturing
func (g *LinuxGrabber) Screen() *Screen {
	return &g.screen
}

// Fps returns the frames per sec. we're capturing
func (g *LinuxGrabber) Fps() int {
	return g.fps
}

func (g *LinuxGrabber) Bitrate() int {
	return g.bitrate
}

func (g *LinuxGrabber) SendEvent(msg []byte) {
	eventType := int(msg[0]);
	var e Event = nil;
	if g.events == nil {
		return;
	}
	if eventType == 1 {
		e = newEvent(msg)
	}

    if e != nil {
    	g.events <- e;
    }

}
// NewProvider returns an X Server-based video provider
func NewProvider() (Service, error) {
	return &LinuxProvider{}, nil
}
