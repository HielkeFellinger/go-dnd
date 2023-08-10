package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/controller"
)

func HandleControllerRoutes(router *gin.Engine) {
	// Home (Page) Routes
	router.GET("/", controller.HomePage)

	// User Routes
	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", controller.LoginPage)
		userRoutes.GET("/register", controller.RegisterPage)
	}
}

func HandleStaticContent(router *gin.Engine) {
	router.Static("/assets", "web/assets")
}

func HandleTemplates(router *gin.Engine) {
	router.LoadHTMLGlob("web/templates/*")
}
