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

type ValidationField struct {
	Name  string
	Error string
}

type ValidationFields []ValidationField

type ValidationErrorStruct struct {
	Error  string
	Fields ValidationFields
}

func ValidationError(w http.ResponseWriter, validation *ValidationFields, e error) {
	js, _ := json.Marshal(ValidationErrorStruct{
		Error:  e.Error(),
		Fields: *validation,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(js)
}
