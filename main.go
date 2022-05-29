package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/konnovK/superchat/docs"
	"github.com/konnovK/superchat/internal/handlers/gateways"
	"github.com/konnovK/superchat/internal/migrations"
	"github.com/konnovK/superchat/internal/utils"
	"github.com/konnovK/superchat/internal/workers"
	httpSwagger "github.com/swaggo/http-swagger"
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

	err = migrations.Migrate(db)
	if err != nil {
		panic(err)
	}
	pool, err := utils.NewRedisPool()
	if err != nil {
		panic(err)
	}
	worker := workers.NewWorker(pool)

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	gw := gateways.NewChatGateway(db, worker)
	gw.InitHandlers(r)

	http.Handle("/", r)
	http.ListenAndServe(":8082", nil)
}
