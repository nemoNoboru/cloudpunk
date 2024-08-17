package handlers

import (
	"cloudpunk/cloud"
	"cloudpunk/serial"
	"fmt"
	"net/http"
	"strings"
)

func HandleAPI(w http.ResponseWriter, req *http.Request) {

	bytes, err := serial.EncodeRequestToJSON(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	result, err := cloud.LuaRun("cloudpunk.serverless.%s", strings.Split(req.URL.Path, "/")[2], bytes)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, result)
}

func HandlePage(w http.ResponseWriter, req *http.Request) {

	bytes, err := serial.EncodeRequestToJSON(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	result, err := cloud.LuaRun("cloudpunk.pages.%s", strings.Split(req.URL.Path, "/")[2], bytes)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, result)
}

func HandleStatic(w http.ResponseWriter, req *http.Request) {
	label := strings.Split(req.URL.Path, "/")[1]

	result := cloud.StorageGet(label)

	w.Write(result)
}
