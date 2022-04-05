package model

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Title       string
	Author      string
	TimeToLive  int
	Tags        []Tag `gorm:"many2many:chats_tags;"`
	Messages    []Message
	LastMessage Message
}

type Message struct {
	gorm.Model
	Author string
	Text   string
	Time   time.Time
	ChatID uint
}

type Tag struct {
	gorm.Model
	Title string
}
