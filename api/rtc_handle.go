package api

import (
	"encoding/json"
	"net/http"

	"github.com/wadahana/ga/display"
	"github.com/wadahana/ga/rtc"
	"github.com/wadahana/memu/log"
)


// MakeHandler returns an HTTP handler for the session service
func MakeApiHandler(webrtc rtc.Service, disp display.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		log.Debugf("/session ->\r\n");
		
		dec := json.NewDecoder(r.Body)
		req := newSessionRequest{}
		
		if err := dec.Decode(&req); err != nil {
			handleError(w, err)
			return
		}
		log.Debugf("offer: %s\r\n", req.Offer);
		peer, err := webrtc.CreateRemoteScreenConnection(req.Screen)
		if err != nil {
			handleError(w, err)
			return
		}
		answer, err := peer.ProcessOffer(req.Offer)

		if err != nil {
			handleError(w, err)
			return
		}
		log.Debugf("answer: %s\r\n", answer);
		payload, err := json.Marshal(newSessionResponse{
			Answer: answer,
		})
		if err != nil {
			handleError(w, err)
			return
		}

		w.Write(payload)
	})

	mux.HandleFunc("/screens", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		log.Debugf("/screens ->");
		screens, err := disp.Screens()
		if err != nil {
			handleError(w, err)
			return
		}

		screensPayload := make([]screenPayload, len(screens))

		for i, s := range screens {
			screensPayload[i].Index = s.Index
			screensPayload[i].Name = s.Name
		}
		payload, err := json.Marshal(screensResponse{
			Screens: screensPayload,
		})
		if err != nil {
			handleError(w, err)
			return
		}

		w.Write(payload)
	})
	log.Info("setup /api handler");
	return mux
}