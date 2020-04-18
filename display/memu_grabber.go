// +build memu

package display

import (
	"image"
	"time"
	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
)


// XVideoProvider implements the rdisplay.Service interface for XServer
type MEmuProvider struct{}

// XScreenGrabber captures video from a X server
type MEmuGrabber struct {
	fps     int
	bitrate int
	screen Screen
	vm     *memu.MEmulator
	frames chan *image.RGBA
	stop   chan struct{}
	events chan Event
	closed bool;
}


// CreateScreenGrabber Creates an screen capturer for the X server
func (*MEmuProvider) CreateGrabber(screen Screen, fps int, bitrate int) (Grabber, error) {
	vm, err := memu.GetEmulator(screen.Name)
	if err != nil {
		return nil, err;
	}
	return &MEmuGrabber{
		screen:  screen,
		bitrate: bitrate,
		fps:     fps,
		vm:      vm, 
		frames:  make(chan *image.RGBA),
		stop:    make(chan struct{}),
		events:  make(chan Event),
		closed:  false,
	}, nil
}

// Screens Returns the available screens to capture
func (x *MEmuProvider) Screens() ([]Screen, error) {
	i := 0;

	m := memu.GetEmulators();
	numScreens := len(*m);
	screens := make([]Screen, numScreens)
	
	for _, e := range *m {
		screens[i] = Screen{
			Name:   e.GetName(),
			FrameRate: e.GetFrameRate(),
			Bitrate: e.GetBitrate(),
			Index:  e.GetIndex(),
			Bounds: e.GetDisplayBounds(),
		};
	}
	return screens, nil
}

// Frames returns a channel that will receive an image stream
func (g *MEmuGrabber) Frames() <-chan *image.RGBA {
	return g.frames
}

// Start initiates the screen capture loop
func (g *MEmuGrabber) Start() {
	delta := time.Duration(1000/g.fps) * time.Millisecond
	go func() {
		for {
			startedAt := time.Now()
			select {
			case <-g.stop:
				close(g.frames)
				close(g.events)
				g.closed = true;
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

				} else {
					log.Info("Unknown Event");
				}
			default:
				img, err := g.vm.CaptureVideo();
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
func (g *MEmuGrabber) Stop() {
	close(g.stop)
}

// Screen returns a pointer to the screen we're capturing
func (g *MEmuGrabber) Screen() *Screen {
	return &g.screen
}

// Fps returns the frames per sec. we're capturing
func (g *MEmuGrabber) Fps() int {
	return g.fps
}

func (g *MEmuGrabber) Bitrate() int {
	return g.bitrate
}
// Send Mouse/Keyboard Event to MEmu
func (g *MEmuGrabber) SendEvent(msg []byte) {
	eventType := int(msg[0]);
	var e Event = nil;
	if g.events == nil && !g.closed {
		return;
	}
	if eventType == 1 {
		e = newMouseEvent(msg)
	} 
    if e != nil {
    	g.events <- e;
    }
}

func (g *MEmuGrabber) onMouseEvent(mouseEvent *MouseEvent) {
	if mouseEvent.getButtonType() == LeftButton {
		ev := memu.NewMouseEvent(mouseEvent.mouseType, mouseEvent.x, mouseEvent.y);
		g.vm.SendEvent(ev);
	}
}
// NewProvider returns an MEmu video provider
func NewProvider() (Service, error) {
	return &MEmuProvider{}, nil
}

/*
func init() {
	log.Printf("startRDP MEmu_1");
	memu.StartRDP("MEmu_1", 1)
}
*/