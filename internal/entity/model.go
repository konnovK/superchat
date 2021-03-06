package entity

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Title    string `gorm:"unique"`
	Creator  string
	TTL      int
	Tags     []Tag `gorm:"many2many:chats_tags;"`
	Messages []Message
}

func (c *Chat) ToChatResponse(chatLastMessages Messages) ChatResponse {
	tags := []string{}
	for _, tag := range c.Tags {
		tags = append(tags, tag.Title)
	}

	response := ChatResponse{
		ID:          c.ID,
		Title:       c.Title,
		Creator:     c.Creator,
		TTL:         c.TTL,
		Tags:        tags,
		LastMessage: chatLastMessages.ToMessageResponse(),
	}
	return response
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

func (m *Message) ToMessageResponse() MessageResponse {
	return MessageResponse{
		Sender:  m.Sender,
		Message: m.Message,
		SentAt:  m.CreatedAt,
	}
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
	Title string `gorm:"unique"`
}
type Tags []Tag
