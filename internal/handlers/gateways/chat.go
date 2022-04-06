package gateways

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/usecase"

	"gorm.io/gorm"
)

type ChatGateway struct {
	ChatDTO usecase.ChatDTO
}

func NewChatGateway(db *gorm.DB) *ChatGateway {
	return &ChatGateway{
		ChatDTO: usecase.NewChatContent(db),
	}
}

func (g *ChatGateway) InitHandlers(r *mux.Router) {
	r.HandleFunc("/chat/active", g.GetActiveChats).Methods("GET")
}

// GetActiveChats godoc
// @Summary Get active chats
// @Description Get all active chats at the moment
// @Tags chat
// @Produce json
// @Success 200 {object} entity.GetActiveChatsResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/active [get]
func (g *ChatGateway) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	chats := entity.GetActiveChatsResponse{}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

// GetMessagesByChatId godoc
// @Summary Get all messages in chat
// @Description Get all messages in chat by chat id
// @Tags chat
// @Produce json
// @Param        id   path      int  true  "Chat ID"
// @Success 200 {object} entity.GetMessagesByChatIdResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/{id}/message [get]
func (g *ChatGateway) GetMessagesByChatId(w http.ResponseWriter, r *http.Request) {
	messages := entity.GetMessagesByChatIdResponse{}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// CreateChat godoc
// @Summary Create new chat
// @Description Create new chat
// @Tags chat
// @Accept json
// @Produce json
// @Param chat body entity.AcceptedChat true "Create chat"
// @Success 200 {object} entity.ChatResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat [post]
func (g *ChatGateway) CreateChat(w http.ResponseWriter, r *http.Request) {
	chat := entity.ChatResponse{}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// SendMessage godoc
// @Summary Send message to chat
// @Description Send message to chat by chat id
// @Tags chat
// @Produce json
// @Param        id   path      int  true  "Chat ID"
// @Param message body entity.AcceptedMessage true "Send message"
// @Success 200 {object} entity.MessageResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/{id}/message [post]
func (g *ChatGateway) SendMessage(w http.ResponseWriter, r *http.Request) {
	message := entity.MessageResponse{}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
