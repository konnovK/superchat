package gateways

import (
	"github.com/gorilla/mux"
)

type Gateway interface {
	InitHandlers(r *mux.Router)
}
