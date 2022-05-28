package gateways

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/usecase"
	"github.com/konnovK/superchat/internal/utils"
	"gorm.io/gorm"
)

type SocketGateway struct {
	Upgrader websocket.Upgrader
	Clients  map[int](map[*websocket.Conn]bool) // FIXME: наверное эта штука может занимать очень много места
	ChatDTO  usecase.ChatDTO
}

func NewSocketGateway(db *gorm.DB) Gateway {
	return &SocketGateway{
		Upgrader: websocket.Upgrader{ // FIXME: наверное Upgrader нужно создавать в другом месте
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Clients: make(map[int](map[*websocket.Conn]bool), 0),
		ChatDTO: usecase.NewChatDTO(db), // чтобы отправлять сообщения
	}
}

func (g *SocketGateway) InitHandlers(r *mux.Router) {
	r.HandleFunc("/chat/{id}", g.ChatSocket)
}

func (g *SocketGateway) WriteMessage(chatId int, message []byte) {
	for conn := range g.Clients[chatId] {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (g *SocketGateway) ChatSocket(w http.ResponseWriter, r *http.Request) {
	idFromUrl := mux.Vars(r)["id"]
	chatId, err := strconv.Atoi(idFromUrl)
	if err != nil {
		utils.JSONError(w, err, http.StatusNotFound)
		return
	}

	connection, err := g.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.JSONError(w, err, http.StatusInternalServerError)
		return
	}
	defer connection.Close()

	// Запоминаем, что чел connection подключился к чату chatID
	if g.Clients[chatId] == nil {
		g.Clients[chatId] = make(map[*websocket.Conn]bool, 0)
	}
	g.Clients[chatId][connection] = true
	// если чел connection отключился от чата chatID, то удаляем его
	defer delete(g.Clients[chatId], connection)

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		// Кто-то отправил сообщение
		// FIXME: надо будет куда-нибудь вынести это
		sendMessageRequest := entity.SendMessageRequest{}
		json.NewDecoder(strings.NewReader(string(message))).Decode(&sendMessageRequest)
		validation, err := sendMessageRequest.Validate()
		if err != nil {
			msg, _ := json.Marshal(validation)
			connection.WriteMessage(websocket.TextMessage, msg)
			continue
		}

		messageResponse, err := g.ChatDTO.SendMessage(chatId, sendMessageRequest)
		if err != nil {
			connection.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}

		msg, _ := json.Marshal(messageResponse)
		g.WriteMessage(chatId, msg)
	}
}
