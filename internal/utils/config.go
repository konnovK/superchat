package utils

import (
	"os"
)

type Config struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func NewConfig() *Config {
	dbUser := os.Getenv("SuperchatDbUser")
	dbPassword := os.Getenv("SuperchatDbPassword")
	dbHost := os.Getenv("SuperchatDbHost")
	dbPort := os.Getenv("SuperchatDbPort")
	dbName := os.Getenv("SuperchatDbName")

	return &Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
	}
}
