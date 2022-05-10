package entity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidatedField struct {
	Name  string
	Error string
}

type ValidatedFields []ValidatedField

func (fields ValidatedFields) ValidationError(w http.ResponseWriter) error {
	if len(fields) == 0 {
		return nil
	}

	js, _ := json.Marshal(struct {
		Error  string          `json:"error"`
		Fields ValidatedFields `json:"fields"`
	}{
		Error:  "validation",
		Fields: fields,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(js)

	return fmt.Errorf("validation error")
}

func (ccr *CreateChatRequest) Validate(w http.ResponseWriter) error {
	fields := ValidatedFields{}

	if ccr.Title == "" {
		fields = append(fields, ValidatedField{
			Name:  "title",
			Error: "field shouldn't be empty",
		})
	}

	if ccr.Creator == "" {
		fields = append(fields, ValidatedField{
			Name:  "creator",
			Error: "field shouldn't be empty",
		})
	}

	if ccr.TTL == 0 {
		fields = append(fields, ValidatedField{
			Name:  "ttl",
			Error: "field shouldn't be empty",
		})
	}

	return fields.ValidationError(w)
}
