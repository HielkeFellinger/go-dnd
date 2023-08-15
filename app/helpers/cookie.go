package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

const AuthCookieName = "Session"

type AuthCookieContent struct {
	ID        uint
	ExpiresAt int64
}

func SetAuthJWTCookie(content AuthCookieContent, c *gin.Context) error {
	// Override / Set default expiration
	content.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	// Generate a JWT token and Sign
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"ID":        content.ID,
		"ExpiresAt": content.ExpiresAt,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err == nil {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(AuthCookieName, tokenString, 3600*24, "", "", false, false)
	}

	return err
}

func ResetCookie(name string, c *gin.Context) {
	c.SetCookie(name, "", -1, "", "", false, false)
}
