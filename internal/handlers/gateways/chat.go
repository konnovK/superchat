package gateways

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/usecase"
	"github.com/konnovK/superchat/internal/utils"
	"github.com/konnovK/superchat/internal/workers"

	"gorm.io/gorm"
)

type ChatGateway struct {
	ChatDTO usecase.ChatDTO
	Worker  *workers.Worker
}

func NewChatGateway(db *gorm.DB, worker *workers.Worker) *ChatGateway {
	return &ChatGateway{
		ChatDTO: usecase.NewChatDTO(db),
		Worker:  worker,
	}
}

func (g *ChatGateway) InitHandlers(r *mux.Router) {
	r.HandleFunc("/chat/active", g.GetActiveChats).Methods("GET")

	r.HandleFunc("/chat", g.CreateChat).Methods("POST")

	r.HandleFunc("/chat/{id}/message", g.GetMessagesByChatId).Methods("GET")

	r.HandleFunc("/chat/{id}/message", g.SendMessage).Methods("POST")
}

// GetActiveChats godoc
// @Summary Get active chats
// @Description Get all active chats at the moment
// @Tags chat
// @Produce json
// @Success 200 {object} entity.GetActiveChatsResponse
// @Failure 400
// @Failure 404
// @Failure 500 {object} utils.JSONErrorStruct
// @Router /chat/active [get]
func (g *ChatGateway) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	chats, err := g.ChatDTO.GetActiveChats()
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}
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
// @Failure 404 {object} utils.JSONErrorStruct
// @Failure 500 {object} utils.JSONErrorStruct
// @Router /chat/{id}/message [get]
func (g *ChatGateway) GetMessagesByChatId(w http.ResponseWriter, r *http.Request) {
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	messages, err := g.ChatDTO.GetMessagesByChatId(chatId)
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// CreateChat godoc
// @Summary Create new chat
// @Description Create new chat
// @Tags chat
// @Accept json
// @Produce json
// @Param chat body entity.CreateChatRequest true "Create chat"
// @Success 200 {object} entity.ChatResponse
// @Failure 400 {object} utils.ValidationErrorStruct
// @Failure 404
// @Failure 500 {object} utils.JSONErrorStruct
// @Router /chat [post]
func (g *ChatGateway) CreateChat(w http.ResponseWriter, r *http.Request) {
	createChatRequest := entity.CreateChatRequest{}
	json.NewDecoder(r.Body).Decode(&createChatRequest)
	validation, err := createChatRequest.Validate()
	if err != nil {
		utils.ValidationError(w, validation, err)
		return
	}

	chatResponce, err := g.ChatDTO.CreateChat(createChatRequest)
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	g.Worker.SetJobTimer(int64(chatResponce.TTL), int(chatResponce.ID))

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponce)
}

// SendMessage godoc
// @Summary Send message to chat
// @Description Send message to chat by chat id
// @Tags chat
// @Produce json
// @Param        id   path      int  true  "Chat ID"
// @Param message body entity.SendMessageRequest true "Send message"
// @Success 200 {object} entity.MessageResponse
// @Failure 400 {object} utils.ValidationErrorStruct
// @Failure 404 {object} utils.JSONErrorStruct
// @Failure 500 {object} utils.JSONErrorStruct
// @Router /chat/{id}/message [post]
func (g *ChatGateway) SendMessage(w http.ResponseWriter, r *http.Request) {
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusNotFound)
		return
	}

	sendMessageRequest := entity.SendMessageRequest{}
	json.NewDecoder(r.Body).Decode(&sendMessageRequest)
	validation, err := sendMessageRequest.Validate()
	if err != nil {
		utils.ValidationError(w, validation, err)
		return
	}

	messages, err := g.ChatDTO.SendMessage(chatId, sendMessageRequest)
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
