package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/controller"
	"github.com/hielkefellinger/go-dnd/app/middelware"
)

func HandleControllerRoutes(router *gin.Engine) {
	// Home (Page) Routes
	router.GET("/", middelware.OptionalAuth, controller.HomePage)

	// User Routes
	userRoutes := router.Group("/u")
	{
		userRoutes.GET("/login", middelware.OptionalAuth, controller.LoginPage)
		userRoutes.POST("/login", middelware.OptionalAuth, controller.Login)
		userRoutes.GET("/logout", controller.Logout)
		userRoutes.POST("/logout", controller.Logout)
		userRoutes.GET("/register", middelware.OptionalAuth, controller.RegisterPage)
		userRoutes.POST("/register", middelware.OptionalAuth, controller.Register)
	}

	campaignRoutes := router.Group("/campaign")
	{
		campaignRoutes.GET("/select", middelware.RequireAuth, controller.CampaignSelectPage)
	}
}

func HandleStaticContent(router *gin.Engine) {
	router.Static("/assets", "web/assets")
}

func HandleTemplates(router *gin.Engine) {
	router.LoadHTMLGlob("web/templates/*")
}
