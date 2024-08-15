package handlers

import (
	"cloudpunk/cloud"
	"fmt"
	"net/http"
	"strings"
)

func HandleAPI(w http.ResponseWriter, req *http.Request) {

	result, err := cloud.LuaRun(strings.Split(req.URL.Path, "/")[2])
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, result)
}

func HandleStatic(w http.ResponseWriter, req *http.Request) {
	label := strings.Split(req.URL.Path, "/")[2]

	result := cloud.StorageGet(label)

	w.Write(result)
}
