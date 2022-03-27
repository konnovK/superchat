package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", nil)
	r.HandleFunc("/products", nil)
	r.HandleFunc("/articles", nil)
	http.Handle("/", r)
}
