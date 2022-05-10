package usecase

import (
	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/repository"
	"gorm.io/gorm"
)

type ChatDTO interface {
	GetActiveChats() (entity.GetActiveChatsResponse, error)

	CreateChat(entity.CreateChatRequest) (entity.CreateChatResponse, error)
}
type ChatContent struct {
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository
}

func NewChatContent(db *gorm.DB) ChatDTO {
	return &ChatContent{
		chatRepository:    repository.NewChatRepository(db),
		messageRepository: repository.NewMessageRepository(db),
	}
}

func (c *ChatContent) GetActiveChats() (entity.GetActiveChatsResponse, error) {
	chats, err := c.chatRepository.FindAll()
	if err != nil {
		return nil, err
	}
	messages := []entity.Message{} // From messages repository
	chatResponse := chats.ToChatResponse(messages)
	return chatResponse, nil
}

func (c *ChatContent) CreateChat(request entity.CreateChatRequest) (entity.CreateChatResponse, error) {
	chat := request.ToChat()

	// TODO: создать теги, если их еще нет ()
	err := c.chatRepository.Create(chat)
	if err != nil {
		return entity.CreateChatResponse{}, err
	}

	response := chat.ToCreateChatResponce()

	return response, nil
}
