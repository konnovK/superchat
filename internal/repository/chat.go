package repository

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/gorm"
)

type ChatRepository interface {
	Find(conditions *entity.Chat) (entity.Chats, error)
	FindAll() (entity.Chats, error)
	Create(target *entity.Chat) error
	Update(conditions *entity.Chat, target *entity.Chat) error
	Delete(target *entity.Chat) error
}

type Chat struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *Chat {
	return &Chat{
		db: db,
	}
}

func (c *Chat) Find(conditions *entity.Chat) (entity.Chats, error) {
	chats := []entity.Chat{}

	queryResult := c.db.Where(conditions).Find(&chats)
	if queryResult.Error != nil {
		return chats, queryResult.Error
	}

	return chats, nil
}

func (c *Chat) FindAll() (entity.Chats, error) {
	return c.Find(&entity.Chat{})
}

func (c *Chat) Create(target *entity.Chat) error {
	queryResult := c.db.Omit("ID").Create(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}
	return nil
}

func (c *Chat) Update(conditions *entity.Chat, target *entity.Chat) error {
	queryResult := c.db.Where(&conditions).Updates(target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}

func (c *Chat) Delete(target *entity.Chat) error {
	queryResult := c.db.Delete(target)

	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}
