package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/konnovK/superchat/docs"
	"github.com/konnovK/superchat/internal/handlers/gateways"
	"github.com/konnovK/superchat/internal/migrations"
	"github.com/konnovK/superchat/internal/utils"
	"github.com/konnovK/superchat/internal/workers"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/rs/cors"
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
	log.Println("now listening on port 8082")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowOriginFunc: func(origin string) bool {return true},
		AllowCredentials: true,
	})
	h := c.Handler(r)
	http.ListenAndServe(":8082", h)
}
