package api

import (
	"net/http"
	"encoding/json"
	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
)

type newSessionRequest struct {
	Offer  string  `json:"offer"`
	Screen string  `json:"screen"`
}

type newSessionResponse struct {
	Answer string `json:"answer"`
}

type screenPayload struct {
	Index int `json:"index"`
	Name string `json:"name"`
}

type screensResponse struct {
	Screens []screenPayload `json:"screens"`
}

type VMListContent struct {
	Nums  int             `json:"nums"`
	List []memu.MEmuInfo  `json:"vms"`
}

type VMCreateContent struct {
	Index int      `json:"index"`
	Name  string   `json:"name"`
}

type ResponseWrapper struct {
	ReturnCode int         `json:"returnCode"`
	RetureMsg  string      `json:"returnMsg"`
	Type       int         `json:"type"`
	Content    interface{} `json:"content"`
}

func handleError(w http.ResponseWriter, err error) {
	log.Errorf("http handler error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func MakeResponse(returnCode int, returnMsg string, content interface{}) ([]byte, error) {
	var resp ResponseWrapper
	resp.ReturnCode = returnCode;
	resp.RetureMsg  = returnMsg;
	resp.Type = 0;
	if content == nil {
		resp.Content = "" 
	} else {
		resp.Content = content
	}
	return json.Marshal(resp)
}

func MakeResponseWithError(err *memu.MEmuError, content interface{}) ([]byte, error) {
	var resp ResponseWrapper
	if err == nil {
		err = memu.ErrorSuccess
	}
	resp.ReturnCode = err.GetCode();
	resp.RetureMsg  = err.Error();
	resp.Type = 0;
	if content == nil {
		resp.Content = "" 
	} else {
		resp.Content = content
	}
	return json.Marshal(resp)
}

