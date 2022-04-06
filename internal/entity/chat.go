package entity

import "time"

type GetActiveChatsResponse []ChatResponse

type ChatResponse struct {
	ID          uint
	Title       string
	Author      string
	TTL         int
	Tags        []string
	LastMessage []MessageResponse
}

type MessageResponse struct {
	Sender  string
	Message string
	Time    time.Time
}

type GetMessagesByChatIdResponse []MessageResponse

type AcceptedChat struct {
	Title  string
	Author string
	TTL    int
	Tags   []string
}

type AcceptedMessage struct {
	Message string
}
