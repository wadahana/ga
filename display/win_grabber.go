// +build win

package display

import (
	"image"
	"time"
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/wadahana/wa"
	"github.com/wadahana/memu/log"
)

// WinProvider implements the rdisplay.Service interface for XServer
type WinProvider struct{}

// WinGrabber captures video from a X server
type WinGrabber struct {
	fps     int
	bitrate int
	screen  Screen
	frames  chan *image.RGBA
	stop    chan struct {}
	events  chan Event
	
}

// CreateGrabber Creates an screen capturer for the X server
func (*WinProvider) CreateGrabber(screen Screen, fps int, bitrate int) (Grabber, error) {
	return &WinGrabber {
		screen:  screen,
		fps:     fps,
		bitrate: bitrate,
		frames:  make(chan *image.RGBA),
		stop:    make(chan struct{}),
		events:  make(chan Event),
	}, nil
}

// Screens Returns the available screens to capture
func (x *WinProvider) Screens() ([]Screen, error) {
	numScreens := screenshot.NumActiveDisplays()
	log.Debugf("num of displays: %d", numScreens);
	screens := make([]Screen, numScreens)
	for i := 0; i < numScreens; i++ {
		screens[i] = Screen{
			Index:  i,
			Bitrate: 2048,
			FrameRate: 25,
			Name: fmt.Sprintf("Screen_%d", i),
			Bounds: screenshot.GetDisplayBounds(i),
		}
	}
	return screens, nil
}

// Frames returns a channel that will receive an image stream
func (g *WinGrabber) Frames() <-chan *image.RGBA {
	return g.frames
}

// Start initiates the screen capture loop
func (g *WinGrabber) Start() {
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
					/*
					log.Infof("MouseEvent: %d,%d,%d, %0.4f, %0.4f\n", 
							mouseEvent.getEventType(), 
							mouseEvent.getMouseType(),
							mouseEvent.getButtonType(),
							mouseEvent.x, 
							mouseEvent.y);
					*/
					g.onMouseEvent(mouseEvent)

				} else if e.getEventType() == 2 {
					kbEvent := e.(*KeyboardEvent);
					//log.Infof("KeyboardEvent: %02x, %v, %d", kbEvent.getEventType(), kbEvent.getPress(), kbEvent.getKeyCode());
					g.onKeyboardEvent(kbEvent)
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
func (g *WinGrabber) Stop() {
	close(g.stop)
}

// Screen returns a pointer to the screen we're capturing
func (g *WinGrabber) Screen() *Screen {
	return &g.screen
}

// Fps returns the frames per sec. we're capturing
func (g *WinGrabber) Fps() int {
	return g.fps
}

func (g *WinGrabber) Bitrate() int {
	return g.bitrate
}

func (g *WinGrabber) SendEvent(msg []byte) {
	eventType := int(msg[0]);
	var e Event = nil;
	if g.events == nil {
		return;
	}
	if eventType == 1 {
		e = newMouseEvent(msg)
	} else if eventType == 2 {
		e = newKeyboardEvent(msg)
	}
    if e != nil {
    	g.events <- e;
    }
}

func (g *WinGrabber) onMouseEvent(ev *MouseEvent) {
	if ev.getEventType() != 1 {
		return;
	}
	fx, fy := ev.getPos();
	x := int32(fx * float32(g.screen.Bounds.Dx()))
	y := int32(fy * float32(g.screen.Bounds.Dy()))
	if ev.getMouseType() == MouseMove {
		wa.SetMousePos(x, y);
	} else if ev.getMouseType() == MouseDown {
		wa.SetMouseKey(ev.getButtonType(), true);
	} else if ev.getMouseType() == MouseUp {
		wa.SetMouseKey(ev.getButtonType(), false);
	} else if ev.getMouseType() == MouseWheel {

	}
}

func (g *WinGrabber) onKeyboardEvent(ev *KeyboardEvent) {
	if ev.getEventType() != 2 {
		return
	}
	log.Debugf("press(%v), keyCode(%d)", ev.getPress(), ev.getKeyCode());
	wa.SetKeyCode(ev.getPress(), uint16(ev.getKeyCode()))
}
// NewProvider returns an X Server-based video provider
func NewProvider() (Service, error) {
	return &WinProvider{}, nil
}
