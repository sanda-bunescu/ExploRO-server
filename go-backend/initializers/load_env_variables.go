package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvFiles() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
