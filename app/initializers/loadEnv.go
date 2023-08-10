package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvVariables() {
	log.Println("INIT: Attempting loading .env variables")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("INIT: Failure loading .env File", err)
	}
}
