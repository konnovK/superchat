package usecase

import (
	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/repository"
	"gorm.io/gorm"
)

type ChatDTO interface {
	GetActiveChats() (*entity.GetActiveChatsResponse, error)
}
type ChatContent struct {
	repository repository.ChatRepository
	// TODO: message repository and chat repository
}

func NewChatContent(db *gorm.DB) *ChatContent {
	return &ChatContent{
		repository: repository.NewChatRepository(db),
		// TODO: create message repository and chat repository
	}
}

func (c *ChatContent) GetActiveChats() (*entity.GetActiveChatsResponse, error) {
	chats, err := c.repository.FindAll()
	if err != nil {
		return nil, err
	}
	messages := []entity.Message{} // From messages repository
	chatResponse := chats.ToChatResponse(messages)
	return &chatResponse, nil
}
