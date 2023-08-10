package initializers

import (
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
)

func SyncDB() {
	log.Println("INIT: Attempting Sync Database Schema")

	userErr := DB.AutoMigrate(&models.User{})
	if userErr != nil {
		log.Fatal("INIT: Failure Syncing User Mode Schema", userErr)
	}
}
