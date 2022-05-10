package utils

import (
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

	return db, nil
}
