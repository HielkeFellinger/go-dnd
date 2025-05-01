package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hielkefellinger/go-dnd/app/helpers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"net/http"
	"time"
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

	// Check the validity of the password hash; if no user match, the password is an invalid hash or has a cost of 0
	if cost, err := bcrypt.Cost([]byte(user.Password)); err != nil || cost == 0 {
		time.Sleep(5 * time.Second) // Ensures minimal duration in auth attempt
		templateMap[errMessage], templateMap[errTitle] = "Invalid username and or password", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// bcrypt.CompareHashAndPassword will only take time if hash is valid
	if errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); errBcrypt != nil {
		templateMap[errMessage], templateMap[errTitle] = "Invalid username and or password", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	var authCookieContent = helpers.AuthCookieContent{ID: user.ID}
	if errCookie := helpers.SetAuthJWTCookie(authCookieContent, c); errCookie != nil {
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
	if c.Bind(&user) != nil {
		templateMap[errMessage], templateMap[errTitle] = "Failed to read request", "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}
	user.Name = html.EscapeString(user.Name)

	// Attempt to insert user
	var service = models.UserService{}
	if err := service.InsertUser(&user); err != nil {
		templateMap[errMessage], templateMap[errTitle] = err.Error(), "Error"
		c.HTML(http.StatusBadRequest, template, templateMap)
		return
	}

	// Redirect (After creating a successful user)
	c.Redirect(http.StatusFound, "/u/login")
}

func Logout(c *gin.Context) {
	helpers.ResetCookie(helpers.AuthCookieName, c)
	c.Redirect(http.StatusFound, "/")
}
