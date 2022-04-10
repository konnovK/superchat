package utils

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDbSession(dbHost string, dbUser string, dbPassword string, dbName string, dbPort string) (*gorm.DB, error) {
	dsn := "host=" + dbHost +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" port=" + dbPort +
		" sslmode=disable TimeZone=Europe/Moscow"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = migrate(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&entity.Tag{}, &entity.Chat{}, &entity.Message{})

	return err
}
