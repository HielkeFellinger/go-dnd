package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/initializers"
	"github.com/hielkefellinger/go-dnd/app/routes"
	"log"
	"os"
)

var router *gin.Engine

func init() {
	log.Println("INIT: Starting Initialisation of GO-DND")
	initializers.LoadEnvVariables()
	initializers.LoadDatabase()
	initializers.SyncDB()
	log.Println("INIT: Done. Initialisation Finished")
}

func main() {
	log.Println("MAIN: Creation of Gin.Engine")
	router = gin.Default()

	// Init Routes
	log.Println("MAIN: Loading (Static) Content, Templates and Routes")
	routes.HandleStaticContent(router)
	routes.HandleTemplates(router)
	routes.HandleControllerRoutes(router)

	// Serve Content
	log.Println("MAIN: Starting Gin.Engine")
	log.Fatal(router.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))))
}
