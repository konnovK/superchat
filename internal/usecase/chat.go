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
	tagRepository     repository.TagRepository
}

func NewChatDTO(db *gorm.DB) ChatDTO {
	return &ChatContent{
		chatRepository:    repository.NewChatRepository(db),
		messageRepository: repository.NewMessageRepository(db),
		tagRepository:     repository.NewTagRepository(db),
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

	tags := []entity.Tag{}
	for _, tag := range chat.Tags {
		c.tagRepository.Create(&tag)
		tags = append(tags, tag)
	}
	chat.Tags = tags

	err := c.chatRepository.Create(chat)
	if err != nil {
		return entity.CreateChatResponse{}, err
	}

	response := chat.ToCreateChatResponce()

	return response, nil
}
