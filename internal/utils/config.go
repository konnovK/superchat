package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func NewConfig() (*Config, error) {
	godotenv.Load()

	dbUser := os.Getenv("CHAT_DB_USER")
	dbPassword := os.Getenv("CHAT_DB_PASSWORD")
	dbHost := os.Getenv("CHAT_DB_HOST")
	dbPort := os.Getenv("CHAT_DB_PORT")
	dbName := os.Getenv("CHAT_DB_NAME")

	return &Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbName:     dbName,
	}, nil
}
