package api

import (
	"net/http"
	"strings"
	"strconv"

	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
)


func MakeVMHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", 	createHandler)
	mux.HandleFunc("/remove", 	removeHandler)
	mux.HandleFunc("/list", 	listHandler)
	mux.HandleFunc("/clone", 	cloneHandler)
	mux.HandleFunc("/start", 	startVMHandler)
	mux.HandleFunc("/stop",  	stopVMHandler)
	mux.HandleFunc("/reboot", 	rebootHandler)
	mux.HandleFunc("/sendkey", 	sendkeyHandler)
	mux.HandleFunc("/shake", 	shakeHandler)
	mux.HandleFunc("/rename",   renameHandler)
	log.Info("setup /vm handler");
	return mux
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var err *memu.MEmuError = nil
	params := r.URL.Query();
	version := params.Get("version")
	ver := memu.AndroidRomUnsupport; 
	if strings.EqualFold(version, "v4.4") {
		ver = memu.AndroidRomV44
	} else if strings.EqualFold(version, "v5.1") {
		ver = memu.AndroidRomV51
	}  else if strings.EqualFold(version, "v7.1") {
		ver = memu.AndroidRomV71
	} else {
		err = memu.ErrorAndroidVersionNotSupport;
	}
	var payload []byte;
	var _err error = nil;
	if err == nil {
		var content VMCreateContent;
		content.Index, content.Name, err = memu.Cmd.Create(ver)
		payload, _err = MakeResponseWithError(err, content)
	}  else {
		payload, _err = MakeResponseWithError(err, nil)
	}
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)	

}

func removeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err *memu.MEmuError  = nil
	params := r.URL.Query();
	name    := params.Get("name")
	strId   := params.Get("index")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.RemoveById(id)
	} else if len(name) > 0 {
		err = memu.Cmd.RemoveByName(name)
	} else {
		err = memu.ErrorInvalidArgument
	}

	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func renameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var err *memu.MEmuError  = nil
	params := r.URL.Query();
	name    := params.Get("name")
	strId   := params.Get("index")
	newName := params.Get("newName")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.RenameById(id, newName)
	} else if len(name) > 0 {
		err = memu.Cmd.RenameByName(name, newName)
	} else {
		err = memu.ErrorInvalidArgument
	}

	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)

}

func listHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var content interface{} = nil
	var payload []byte;

	list, err := memu.Cmd.List(false);
	if err == nil && len(*list) > 0 {
		listContent := VMListContent{
			Nums:  len(*list),
			List:  *list,
		}
		content = listContent	
	} 
	payload, _err := MakeResponseWithError(err, content)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	payload, _err := MakeResponseWithError(memu.ErrorSuccess, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func startVMHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query();

	var err *memu.MEmuError  = nil
	name    := params.Get("name")
	strId   := params.Get("index")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		log.Debugf("index: %d", id);
		err = memu.Cmd.StartById(id)
	} else if len(name) > 0 {
		log.Debugf("name: %s", name);
		err = memu.Cmd.StartByName(name)
	} else {
		err = memu.ErrorInvalidArgument
	}

	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func stopVMHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query();

	var err *memu.MEmuError  = nil
	name    := params.Get("name")
	strId   := params.Get("index")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.StopById(id)
	} else if len(name) > 0 {
		err = memu.Cmd.StopByName(name)
	} else {
		err = memu.ErrorInvalidArgument
	}

	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func rebootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query();

	var err *memu.MEmuError  = nil
	name    := params.Get("name")
	strId   := params.Get("index")

	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.RebootById(id)
	} else if len(name) > 0 {
		err = memu.Cmd.RebootByName(name)
	} else {
		err = memu.ErrorInvalidArgument
	}
	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func sendkeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query();

	var err *memu.MEmuError  = nil
	name    := params.Get("name")
	strId   := params.Get("index")
	key   	:= params.Get("key")
	id, e := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.SendKeyById(id, key)
	} else if len(name) > 0 {
		err = memu.Cmd.SendKeyByName(name, key)
	} else {
		err = memu.ErrorInvalidArgument
	}
	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}

func shakeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	params := r.URL.Query();

	var err *memu.MEmuError  = nil
	name    := params.Get("name")
	strId   := params.Get("index")
	id, e   := strconv.Atoi(strId)
	if e == nil && id >= 0 && id < 40 {
		err = memu.Cmd.ShakeById(id)
	} else if len(name) > 0 {
		err = memu.Cmd.ShakeByName(name)
	} else {
		err = memu.ErrorInvalidArgument
	}
	payload, _err := MakeResponseWithError(err, nil)
	if _err != nil {
		handleError(w, _err)
		return
	}
	w.Write(payload)
}
