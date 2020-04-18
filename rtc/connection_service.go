package rtc

import (
	"fmt"
	"strings"
	"github.com/wadahana/ga/display"
	"github.com/wadahana/ga/encoders"
)

// RemoteScreenService is our implementation of the rtc.Service
type RemoteScreenService struct {
	videoService    display.Service
	encodingService encoders.Service
}

// NewRemoteScreenService creates a new instances of RemoteScreenService
func NewRemoteScreenService(video display.Service, enc encoders.Service) Service {
	return &RemoteScreenService{
		videoService:    video,
		encodingService: enc,
	}
}

func hasElement(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// CreateRemoteScreenConnection creates and configures a new peer connection
// that will stream the selected screen
func (svc *RemoteScreenService) CreateRemoteScreenConnection(name string) (RemoteScreenConnection, error) {
	screens, err := svc.videoService.Screens()
	if err != nil {
		return nil, err
	}
	var screen *display.Screen = nil;
	for _, s := range screens {
		if strings.EqualFold(s.Name, name) {
			screen = &s;
		} 
	}

	if screen == nil {
		return nil, fmt.Errorf("No available screens")
	}
	screenGrabber, err := svc.videoService.CreateGrabber(*screen, screen.FrameRate, screen.Bitrate)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	if len(screens) == 0 {
		return nil, fmt.Errorf("No available screens")
	}

	rtcPeer := newRemoteScreenPeerConn(screenGrabber, svc.encodingService)
	return rtcPeer, nil
}
