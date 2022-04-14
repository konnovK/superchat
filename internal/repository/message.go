package repository

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Find(conditions *entity.Message) (entity.Messages, error)
	FindAll() (entity.Messages, error)
}

type Message struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *Message {
	return &Message{
		db: db,
	}
}

func (m *Message) Find(conditions *entity.Message) (entity.Messages, error) {
	message := []entity.Message{}

	queryResult := m.db.Where(conditions).Find(&message)
	if queryResult.Error != nil {
		return message, queryResult.Error
	}

	return message, nil
}

func (m *Message) FindAll() (entity.Messages, error) {
	return m.Find(&entity.Message{})
}
