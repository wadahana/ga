package api

import (
	"encoding/json"
	"net/http"

	"github.com/wadahana/ga/display"
	"github.com/wadahana/ga/rtc"
	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
)


// MakeHandler returns an HTTP handler for the session service
func MakeApiHandler(webrtc rtc.Service, disp display.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		
		var peer rtc.RemoteScreenConnection;
		var payload []byte;
		var answer string;
		var err error = nil;
		var _err *memu.MEmuError = nil;

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		log.Debugf("/session ->\r\n");
		
		dec := json.NewDecoder(r.Body)
		req := newSessionRequest{}
		
		err = dec.Decode(&req);

		log.Debugf("offer: %s\r\n", req.Offer);
		if err == nil {
			peer, err = webrtc.CreateRemoteScreenConnection(req.Screen)
		}
		
		answer, err = peer.ProcessOffer(req.Offer)

		log.Debugf("answer: %s\r\n", answer);

		if err != nil {
			_err = memu.NewError(memu.RC_SystemError, err.Error())
			payload, err = MakeResponseWithError(_err, nil);
		} else {
			payload, err = MakeResponseWithError(_err, newSessionResponse{
				Answer: answer,
			})
		}
		
		if err != nil {
			handleError(w, err)
			return
		}
		w.Write(payload)

	})

	mux.HandleFunc("/screens", func(w http.ResponseWriter, r *http.Request) {

		var payload []byte;
		var err error = nil;
		var _err *memu.MEmuError = nil;
		var screensPayload []screenPayload = nil;
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		log.Debugf("/screens ->");

		screens, err := disp.Screens()
		if err == nil {

			screensPayload := make([]screenPayload, len(screens))

			for i, s := range screens {
				screensPayload[i].Index = s.Index
				screensPayload[i].Name = s.Name
			}
		}

		if err != nil {
			_err = memu.NewError(memu.RC_SystemError, err.Error())
			payload, err = MakeResponseWithError(_err, nil);
		} else {
			payload, err = MakeResponseWithError(_err, screensResponse{
				Screens: screensPayload,
			})
		}
		
		if err != nil {
			handleError(w, err)
			return
		}
		w.Write(payload)
	})
	log.Info("setup /api handler");
	return mux
}