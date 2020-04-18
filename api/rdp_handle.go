package api

import(
	"net/http"
	"strconv"
	"github.com/wadahana/ga/core"
	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
)


func MakeRdpHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/start", startRdpHandler)
	mux.HandleFunc("/stop",  stopRdpHandler)
	log.Info("setup /rdp handler");
	return mux
}

func startRdpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var err *memu.MEmuError = nil;
	var vmInfo *memu.MEmuInfo = nil

	params := r.URL.Query();
	name         := params.Get("name")
	strId        := params.Get("index")
	strFps       := params.Get("fps")
	strBitrate   := params.Get("bitrate")

	fps, e := strconv.Atoi(strFps)
	if e != nil || fps < 5 || fps > 60 {
		fps = core.Setting.Rtc.FrameRate
	} 
	bitrate, e := strconv.Atoi(strBitrate)
	if e != nil || bitrate < 64 || bitrate > 4096 {
		bitrate = core.Setting.Rtc.Bitrate
	}
	

	id, e := strconv.Atoi(strId)
	log.Debugf("e: %v, id: %d, strId:%s", e, id, strId);
	if e == nil && id >= 0 && id < 40 {
		vmInfo, err = memu.Cmd.LookupByIndex(id, true)
	} else if len(name) > 0 {
		vmInfo, err = memu.Cmd.LookupByName(name, true)
	} else {
		err = memu.ErrorInvalidArgument
	}
	
	if err == nil && vmInfo == nil {
		err = memu.ErrorEmulatorNotRunning
	} 
	if err == nil {
		log.Debugf("start rdp for vm(%s) fps(%d), bitrate(%dbps)", vmInfo.Name, fps, bitrate)
		//err = rdp.Start(vmInfo.Index, vmInfo.Name, fps, bitrate);
		err = memu.Cmd.StartMiracast(vmInfo.Name)
		if err == nil {
			err = memu.StartRDP(vmInfo.Name, vmInfo.Index, fps, bitrate);
		}
	} else {
		log.Debugf("start rdp fail, vm(%d,%s) not found!", id, name)
	}

	payload, _err := MakeResponseWithError(err, nil)

	if _err != nil {
		handleError(w, _err)
		return
	}

	w.Write(payload)	
}

func stopRdpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err *memu.MEmuError = nil;
	var vmInfo *memu.MEmuInfo = nil
	params := r.URL.Query();
	name    := params.Get("name")
	strId   := params.Get("index")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		vmInfo, err = memu.Cmd.LookupByIndex(id, true)
	} else if len(name) > 0 {
		vmInfo, err = memu.Cmd.LookupByName(name, true)
	} else {
		err = memu.ErrorInvalidArgument
	}

	if err == nil {
		log.Debugf("stop rdp for vm(%s) ", vmInfo.Name)
		memu.StopRDP(vmInfo.Name);
		//err = memu.Cmd.StopMiracast(vmInfo.Name)
	} else {
		log.Debugf("stop rdp fail, vm(%d, %s) not found", vmInfo.Name , vmInfo.Index)
	}
	payload, _err := MakeResponseWithError(err, nil)

	if _err != nil {
		handleError(w, _err)
		return
	}

	w.Write(payload)	
}


