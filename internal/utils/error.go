package utils

import (
	"encoding/json"
	"net/http"
)

type JSONErrorStruct struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, e error, status int) {
	var d JSONErrorStruct
	d.Error = e.Error()
	js, _ := json.Marshal(d)
	http.Error(w, string(js), status)
}
