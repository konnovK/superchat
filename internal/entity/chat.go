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
	Title   string
	Creator string
	TTL     int
	Tags    []string
}

func (ccr CreateChatRequest) ToChat() Chat {
	tags := []Tag{}
	for _, t := range ccr.Tags {
		tags = append(tags, Tag{
			Title: t,
		})
	}

	return Chat{
		Title:   ccr.Title,
		Creator: ccr.Creator,
		TTL:     ccr.TTL,
		Tags: tags,
	}
}

type SendMessageRequest struct {
	Sender  string
	Message string
}

type CreateChatResponse struct {
	ID          uint
	Title       string
	Creator     string
	TTL         int
	Tags        []string
	LastMessages []MessageResponse
}
func (c Chat) ToCreateChatResponce() CreateChatResponse {
	tags := []string{}
	messages := []MessageResponse{}

	for _, t := range c.Tags {
		tags = append(tags, t.Title)
	}

	ccr := CreateChatResponse{
		ID: c.ID,
		Title: c.Title,
		Creator: c.Creator,
		TTL: c.TTL,
		Tags: tags,
		LastMessages: messages,
	}
	return ccr
}