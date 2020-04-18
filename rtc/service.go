package rtc

import (
	"io"
)

type videoStreamer interface {
	start()
	close()
	onMessage(msg []byte)
}

// RemoteScreenConnection Represents a WebRTC connection to a single peer
type RemoteScreenConnection interface {
	io.Closer
	ProcessOffer(offer string) (string, error)
}

// Service WebRTC service
type Service interface {
	CreateRemoteScreenConnection(name string) (RemoteScreenConnection, error)
}
