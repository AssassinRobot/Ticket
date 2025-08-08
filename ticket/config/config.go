package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (error) {
	err:= godotenv.Load(".env")

	return  err
}

func GetDatabaseURL() string {
	return  os.Getenv("DATABASE_URL")
}

func GetNatsURL() string {
	return os.Getenv("NATS_URL")
}