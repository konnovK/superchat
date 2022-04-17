package entity

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Title    string
	Creator  string
	TTL      int
	Tags     []Tag `gorm:"many2many:chats_tags;"`
	Messages []Message
}
type Chats []Chat

func (c Chats) ToChatResponse(chatLastMessages Messages) GetActiveChatsResponse {
	var response GetActiveChatsResponse
	for _, v := range c {
		var tags []string
		for _, tag := range v.Tags {
			tags = append(tags, tag.Title)
		}
		messages := chatLastMessages.ToMessageResponse()
		response = append(response, ChatResponse{
			ID:          v.ID,
			Title:       v.Title,
			TTL:         v.TTL,
			Creator:     v.Creator,
			Tags:        tags,
			LastMessage: messages,
		})
	}
	return response
}

type Message struct {
	gorm.Model
	Sender  string
	Message string
	ChatID  uint
}
type Messages []Message

func (m Messages) ToMessageResponse() []MessageResponse {
	var response []MessageResponse
	for _, message := range m {
		var mr MessageResponse
		mr.Sender = message.Sender
		mr.Message = message.Message
		mr.SentAt = message.CreatedAt
		response = append(response, mr)
	}

	return response
}

type Tag struct {
	gorm.Model
	Title string
}
type Tags []Tag
