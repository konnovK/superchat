package repository

import (
	"github.com/konnovK/superchat/internal/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	Find(conditions model.Chat) (model.Chat, error)
	FindAll() ([]model.Chat, error)
	Create(target model.Chat) error
	Update(conditions model.Chat, target model.Chat) error
}

type Chat struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *Chat {
	return &Chat{
		db: db,
	}
}

func (c *Chat) Find(conditions model.Chat) (model.Chat, error) {
	chat := model.Chat{}

	queryResult := c.db.Where(conditions).First(&chat)
	if queryResult.Error != nil {
		return chat, queryResult.Error
	}

	return chat, nil
}

func (c *Chat) FindAll() ([]model.Chat, error) {
	chats := []model.Chat{}

	queryResult := c.db.Find(&chats)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return chats, nil
}

func (c *Chat) Create(target model.Chat) error {
	queryResult := c.db.Omit("ID").Create(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}
	return nil
}

func (c *Chat) Update(conditions model.Chat, target model.Chat) error {
	queryResult := c.db.Where(&conditions).Updates(target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}
