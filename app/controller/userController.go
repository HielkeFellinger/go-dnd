package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func LoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{"title": "GO-DND Login"},
	)
}

func Login(c *gin.Context) {
	var body struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Login"
	const template = "login.html"

	if c.Bind(&body) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Check if user exists
	var user models.User
	models.DB.First(&user, "name = ?", body.Username)
	if user.ID == 0 {
		templateMap[errMessage], templateMap[errTitle] = "Invalid username and or password", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errBcrypt != nil {
		templateMap[errMessage], templateMap[errTitle] = "Invalid username and or password", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	var authCookieContent = helpers.AuthCookieContent{ID: user.ID}
	errCookie := helpers.SetAuthJWTCookie(authCookieContent, c)
	if errCookie != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed create Cookie", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect
	c.Redirect(http.StatusFound, "/campaign/select")
}

func RegisterPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"register.html",
		gin.H{"title": "GO-DND Register"},
	)
}

func Register(c *gin.Context) {
	templateMap := gin.H{}
	templateMap["title"] = "GO-DND Register"
	const template = "register.html"

	// Parse body to model
	var user models.User
	var service = models.UserService{}
	if c.Bind(&user) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Attempt to insert user
	err := service.InsertUser(&user)
	if err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect (After creating a successful user)
	c.Redirect(http.StatusCreated, "/u/login")
}

func Logout(c *gin.Context) {
	helpers.ResetCookie(helpers.AuthCookieName, c)
	c.Redirect(http.StatusFound, "/")
}
