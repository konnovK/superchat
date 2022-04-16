package repository

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Find(conditions *entity.Message) (entity.Messages, error)
	FindAll() (entity.Messages, error)
	Create(target *entity.Message) error
	Update(conditions *entity.Message, target *entity.Message) error
	Delete(target *entity.Message) error
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

func (m *Message) Create(target *entity.Message) error {
	queryResult := m.db.Omit("ID").Create(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}
	return nil
}

func (m *Message) Update(conditions *entity.Message, target *entity.Message) error {
	queryResult := m.db.Where(&conditions).Updates(target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}

func (m *Message) Delete(target *entity.Message) error {
	queryResult := m.db.Where(&target).Delete(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}
