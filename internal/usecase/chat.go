package usecase

import (
	"github.com/konnovK/superchat/internal/entity"
	"github.com/konnovK/superchat/internal/repository"
	"gorm.io/gorm"
)

type ChatDTO interface {
	GetActiveChats() (entity.GetActiveChatsResponse, error)

	CreateChat(entity.CreateChatRequest) (entity.CreateChatResponse, error)

	GetMessagesByChatId(int) (entity.GetMessagesByChatIdResponse, error)

	SendMessage(int, entity.SendMessageRequest) (entity.MessageResponse, error)
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
	chatsResponse := entity.GetActiveChatsResponse{}
	for _, chat := range chats {
		// lastMessages := []entity.Message{}
		lastMessages, err := c.messageRepository.FindLastMessagesByChatId(chat.ID)
		if err != nil {
			return nil, err
		}
		chatsResponse = append(chatsResponse, chat.ToChatResponse(lastMessages))
	}
	return chatsResponse, nil
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

func (c *ChatContent) GetMessagesByChatId(chatId int) (entity.GetMessagesByChatIdResponse, error) {
	messages, err := c.messageRepository.Find(&entity.Message{ChatID: uint(chatId)})
	if err != nil {
		return nil, err
	}
	messageResponce := messages.ToMessageResponse()
	return messageResponce, nil
}

func (c *ChatContent) SendMessage(chatId int, messageRequest entity.SendMessageRequest) (entity.MessageResponse, error) {
	message := messageRequest.ToMessage(uint(chatId))
	err := c.messageRepository.Create(message)
	if err != nil {
		return entity.MessageResponse{}, err
	}
	return message.ToMessageResponse(), nil
}
