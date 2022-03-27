package gateways

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/konnovK/superchat/internal/usecase"

	"github.com/konnovK/superchat/internal/utils"
)

type ChatGateway struct {
	ChatDTO usecase.ChatDTO
}

func (g *ChatGateway) InitHandlers(r *mux.Router) {
	r.HandleFunc("/chat/active", g.GetActiveChats).Methods("GET")
}

// GetActiveChats godoc
// @Summary get active chats
// @Description get all active chats at the moment
// @Tags chat
// @Produce json
// @Success 200 {array} entity.GetActiveChatsResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/active [get]
func (g *ChatGateway) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	allChats, err := g.ChatDTO.GetAllChats()
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(allChats)
}
