package gateways

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/usecase"
	"github.com/konnovK/superchat/internal/utils"
	"github.com/konnovK/superchat/internal/workers"

	"gorm.io/gorm"
)

type ChatGateway struct {
	ChatDTO usecase.ChatDTO
	Worker  *workers.Worker

	SendedMessage chan entity.MessageResponse
	WebSocket     struct {
		Upgrader websocket.Upgrader
		Clients  map[int](map[*websocket.Conn]time.Time) // FIXME: наверное эта штука может занимать очень много места
	}
}

func NewChatGateway(db *gorm.DB, worker *workers.Worker) *ChatGateway {
	return &ChatGateway{
		ChatDTO:       usecase.NewChatDTO(db),
		Worker:        worker,
		SendedMessage: make(chan entity.MessageResponse, 1),
		WebSocket: struct {
			Upgrader websocket.Upgrader
			Clients  map[int](map[*websocket.Conn]time.Time)
		}{
			Upgrader: websocket.Upgrader{ // FIXME: наверное Upgrader нужно создавать в другом месте
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			Clients: make(map[int](map[*websocket.Conn]time.Time), 0),
		},
	}
}

func (g *ChatGateway) InitHandlers(r *mux.Router) {
	r.HandleFunc("/chat/active", g.GetActiveChats).Methods("GET")

	r.HandleFunc("/chat", g.CreateChat).Methods("POST")

	r.HandleFunc("/chat/{id}/message", g.GetMessagesByChatId).Methods("GET")

	r.HandleFunc("/chat/{id}/message", g.SendMessage).Methods("POST")

	r.HandleFunc("/chat/{id}", g.ChatSocket)

	r.HandleFunc("/tag", g.GetAllTags).Methods("GET")
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
	chats, err := g.ChatDTO.GetActiveChats()
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

// GetAllTags godoc
// @Summary Get all tags
// @Description Get all tags at the moment
// @Tags tag
// @Produce json
// @Success 200 {object} entity.TagsResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /tag [get]
func (g *ChatGateway) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := g.ChatDTO.GetAllTags()
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
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
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	messages, err := g.ChatDTO.GetMessagesByChatId(chatId)
	if err != nil {
		utils.JSONError(w, err, http.StatusBadRequest)
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
// @Failure 400
// @Failure 404
// @Failure 500
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
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chat/{id}/message [post]
func (g *ChatGateway) SendMessage(w http.ResponseWriter, r *http.Request) {
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	sendMessageRequest := entity.SendMessageRequest{}
	json.NewDecoder(r.Body).Decode(&sendMessageRequest)
	validation, err := sendMessageRequest.Validate()
	if err != nil {
		utils.ValidationError(w, validation, err)
		return
	}

	message, err := g.ChatDTO.SendMessage(chatId, sendMessageRequest)
	if err != nil {
		utils.JSONError(w, err, http.StatusBadRequest)
		return
	}

	go func() {
		g.SendedMessage <- message
	}()

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (g *ChatGateway) ChatSocket(w http.ResponseWriter, r *http.Request) {
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusNotFound)
		return
	}

	connection, err := g.WebSocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	defer connection.Close()

	// Запоминаем, что чел connection подключился к чату chatID
	if g.WebSocket.Clients[chatId] == nil {
		g.WebSocket.Clients[chatId] = make(map[*websocket.Conn]time.Time, 0)
	}
	g.WebSocket.Clients[chatId][connection] = time.Now().UTC()
	// если чел connection отключился от чата chatID, то удаляем его
	defer delete(g.WebSocket.Clients[chatId], connection)

	for {
		messageResponses := <-g.SendedMessage

		msg, _ := json.Marshal(messageResponses)
		for conn := range g.WebSocket.Clients[chatId] {

			if g.WebSocket.Clients[chatId][conn].Before(messageResponses.SentAt) {
				conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
