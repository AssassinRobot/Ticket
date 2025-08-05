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
	return os.Getenv("Nats_URL")
}

func GetPort() string {
	return os.Getenv("PORT")
}

func GetENV() string {
	return os.Getenv("ENV")
}


