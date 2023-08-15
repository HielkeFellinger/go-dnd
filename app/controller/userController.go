package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hielkefellinger/go-dnd/app/initializers"
	"github.com/hielkefellinger/go-dnd/app/models"
	"github.com/hielkefellinger/go-dnd/app/util"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
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

	if c.Bind(&body) != nil {
		handeError(c, "login.html", "GO-DND Login", "Failed to read request", "Error")
		return
	}

	// Check if user exists
	var user models.User
	initializers.DB.First(&user, "name = ?", body.Username)
	if user.ID == 0 {
		handeError(c, "login.html", "GO-DND Login", "Invalid username and or password", "Error")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		handeError(c, "login.html", "GO-DND Login", "Invalid username and or password", "Error")
		return
	}

	// Generate a JWT token and Sign
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"ID":        user.ID,
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		handeError(c, "login.html", "GO-DND Login", "Failure while setting auth. token", "Error")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Session", tokenString, 3600*24, "", "", false, false)

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

	if c.Bind(&body) != nil {
		handeError(c, "register.html", "GO-DND Register", "Failed to read request", "Error")
		return
	}

	if body.PasswordCheck != body.Password {
		handeError(c, "register.html", "GO-DND Register", "Passwords do not match", "Error")
		return
	}

	hashByteArray, err := util.HashPassword(body.Password)
	if err != nil {
		handeError(c, "register.html", "GO-DND Register", "Password could not be hashed", "Error")
		return
	}

	// Attempt to create user
	user := models.User{Name: body.Username, Password: string(hashByteArray)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		handeError(c, "register.html", "GO-DND Register", "User could not created", "Error")
		return
	}

	// Redirect
	c.Redirect(http.StatusCreated, "/u/login")
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
