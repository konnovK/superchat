package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/konnovK/superchat/internal/handlers/gateways"
	"github.com/konnovK/superchat/internal/utils"
)

func main() {
	config, err := utils.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := utils.InitDbSession(config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	gw := gateways.NewChatGateway(db)
	gw.InitHandlers(r)

	http.Handle("/", r)
	http.ListenAndServe(":8082", nil)
}
