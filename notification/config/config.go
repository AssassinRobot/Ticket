package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (error) {
	err:= godotenv.Load(".env")

	return  err
}

func GetNatsURL() string {
	return os.Getenv("NATS_URL")
}

func GetENV() string {
	return os.Getenv("ENV")
}

func GetSMTPConfigs() string {
	return os.Getenv("SMTP")
}
