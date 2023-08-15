package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/initializers"
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

	const loginTemplate = "login.html"
	const title = "GO-DND Login"

	if c.Bind(&body) != nil {
		handeError(c, loginTemplate, title, "Failed to read request", "Error")
		return
	}

	// Check if user exists
	var user models.User
	initializers.DB.First(&user, "name = ?", body.Username)
	if user.ID == 0 {
		handeError(c, loginTemplate, title, "Invalid username and or password", "Error")
		return
	}

	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errBcrypt != nil {
		handeError(c, loginTemplate, title, "Invalid username and or password", "Error")
		return
	}

	var authCookieContent = helpers.AuthCookieContent{ID: user.ID}
	errCookie := helpers.SetAuthJWTCookie(authCookieContent, c)
	if errCookie != nil {
		handeError(c, loginTemplate, title, "Failed create Cookie", "Error")
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
	var body struct {
		Username      string `form:"username"`
		Password      string `form:"password"`
		PasswordCheck string `form:"passwordCheck"`
	}

	const registerTemplate = "register.html"
	const title = "GO-DND Register"

	if c.Bind(&body) != nil {
		handeError(c, registerTemplate, title, "Failed to read request", "Error")
		return
	}

	if body.PasswordCheck != body.Password {
		handeError(c, registerTemplate, title, "Passwords do not match", "Error")
		return
	}

	hashByteArray, err := helpers.HashPassword(body.Password)
	if err != nil {
		handeError(c, registerTemplate, title, "Password could not be hashed", "Error")
		return
	}

	// Attempt to create user
	user := models.User{Name: body.Username, Password: string(hashByteArray)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		handeError(c, registerTemplate, title, "User could not created", "Error")
		return
	}

	// Redirect
	c.Redirect(http.StatusCreated, "/u/login")
}

func Logout(c *gin.Context) {
	helpers.ResetCookie(helpers.AuthCookieName, c)
	c.Redirect(http.StatusFound, "/")
}

func handeError(c *gin.Context, template string, title string, errorMessage string, errorTitle string) {
	c.HTML(
		http.StatusBadRequest,
		template,
		gin.H{"title": title,
			"ErrorMessage": errorMessage,
			"ErrorTitle":   errorTitle},
	)
}
