package middelware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hielkefellinger/go-dnd/app/models"
	"log"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	user, err := retrieveUserFromCookie(c)
	if err != nil || user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user", user)
	c.Next()
}

func OptionalAuth(c *gin.Context) {
	user, err := retrieveUserFromCookie(c)
	log.Println(user.ID)
	if err == nil && user.ID != 0 {
		c.Set("user", user)
	}

	c.Next()
}

func retrieveUserFromCookie(c *gin.Context) (models.User, error) {
	var user models.User

	// Get Cookie (contents)
	tokenString, err := c.Cookie("Session")
	if err == nil {
		// Parse tokenString
		token, jwtErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what is expected
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// Send jwtErr if failure in parsing
		if jwtErr != nil {
			return user, jwtErr
		}

		// Validate the cookie content and attempt to retrieve user
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check the expiration.
			if float64(time.Now().Unix()) > claims["ExpiresAt"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Get user, if exists
			models.DB.First(&user, claims["ID"])
		}
	}
	return user, err
}
