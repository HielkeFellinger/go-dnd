package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/initializers"
	"github.com/hielkefellinger/go-dnd/app/routes"
	"log"
	"os"
)

var engine *gin.Engine

func init() {
	log.Println("INIT: Starting Initialisation of GO-DND")
	initializers.LoadEnvVariables()
	initializers.LoadDatabase()
	initializers.SyncDB()
	log.Println("INIT: Done. Initialisation Finished")
}

func main() {
	loadGinEngine()

	// Serve Content
	log.Println("MAIN: Starting Gin.Engine")
	log.Fatal(engine.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))))
}

func loadGinEngine() {
	log.Println("MAIN: Creation of Gin.Engine")
	engine = gin.Default()

	// Load Routes and (static) content
	log.Println("MAIN: Loading (Static) Content, Templates and Routes")
	routes.HandleStaticContent(engine)
	routes.HandleTemplates(engine)
	routes.HandleControllerRoutes(engine)
}
