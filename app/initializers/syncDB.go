package initializers

import (
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
)

func SyncDB() {
	log.Println("INIT: Attempting Sync Database Schema")

	if DB.AutoMigrate(&models.User{}) != nil {
		log.Fatal("INIT: Failure Syncing User Schema")
	}
	if DB.AutoMigrate(&models.Character{}) != nil {
		log.Fatal("INIT: Failure Syncing Character Schema")
	}
	if DB.AutoMigrate(&models.Campaign{}) != nil {
		log.Fatal("INIT: Failure Syncing Campaign Schema")
	}
}
