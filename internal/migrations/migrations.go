package migrations

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Tag{}, &entity.Chat{}, &entity.Message{})

	return err
}
