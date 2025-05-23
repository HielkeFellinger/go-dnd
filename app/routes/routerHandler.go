package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/controller"
	"github.com/hielkefellinger/go-dnd/app/middelware"
	"github.com/hielkefellinger/go-dnd/app/session"
	"net/http"
	"os"
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
		campaignRoutes.GET("/new", middelware.RequireAuth, controller.CampaignNewPage)
		campaignRoutes.POST("/new", middelware.RequireAuth, controller.CampaignNew)
		campaignRoutes.GET("/login/:id", middelware.RequireAuth, controller.CampaignSessionAuthorize)
		campaignRoutes.POST("/login/:id", middelware.RequireAuth, controller.CampaignSessionAuthorize)
		campaignRoutes.GET("/session/:id", middelware.RequireAuth, controller.CampaignSessionPage)
		campaignRoutes.GET("/session/:id/ws", middelware.RequireAuthAndCampaign, session.ServeSessionWS)
	}
}

func HandleStaticContent(router *gin.Engine) {
	router.Static("/assets", "web/assets")
	router.Static("/images", "web/images")

	campaignDataRoutes := router.Group("/campaign_data", middelware.RequireAuthAndCampaignAccess)
	{
		campaignDataRoutes.StaticFS("", http.Dir(os.Getenv("CAMPAIGN_DATA_DIR")))
	}
}

func HandleTemplates(router *gin.Engine) {
	router.LoadHTMLGlob(os.Getenv("CAMPAIGN_WEB_DIR") + "/templates/*")
}
