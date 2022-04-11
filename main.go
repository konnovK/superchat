package main

import (
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/konnovK/superchat/internal/handlers/gateways"
	"github.com/konnovK/superchat/internal/utils"
	"github.com/konnovK/superchat/internal/workers"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	},
}

func main() {
	config, err := utils.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := utils.InitDbSession(config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	if err != nil {
		panic(err)
	}

	worker := workers.NewWorker(redisPool)

	r := mux.NewRouter()

	gw := gateways.NewChatGateway(db, worker)
	gw.InitHandlers(r)

	http.Handle("/", r)
	http.ListenAndServe(":8082", nil)
}
