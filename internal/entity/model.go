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

func (c Chats) ToChatResponse(chatLastMessages []Message) GetActiveChatsResponse {
	var response GetActiveChatsResponse
	for _, v := range c {
		var tags []string
		for _, tag := range v.Tags {
			tags = append(tags, tag.Title)
		}
		messages := []MessageResponse{} // chatLastMessages -> []MessageResponse
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

type Tag struct {
	gorm.Model
	Title string
}
