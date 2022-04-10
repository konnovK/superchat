package entity

import "time"

type GetActiveChatsResponse []ChatResponse

type ChatResponse struct {
	ID          uint
	Title       string
	Creator     string
	TTL         int
	Tags        []string
	LastMessage []MessageResponse
}

type MessageResponse struct {
	Sender  string
	Message string
	SentAt  time.Time
}

type GetMessagesByChatIdResponse []MessageResponse

type CreateChatRequest struct {
	Title  string
	Author string
	TTL    int
	Tags   []string
}

type SendMessageRequest struct {
	Sender  string
	Message string
}
