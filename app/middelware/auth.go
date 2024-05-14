package middelware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hielkefellinger/go-dnd/app/models"
	"net/http"
	"os"
	"strings"
	"time"
)

const loginPageLocation = "/u/login"

func RequireAuth(c *gin.Context) {
	user, err := retrieveUserFromCookie(c)
	if err != nil || user.ID == 0 {
		c.Redirect(http.StatusFound, loginPageLocation)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user", user)
	c.Next()
}

func RequireAuthAndCampaign(c *gin.Context) {
	// Validate User
	user, err := retrieveUserFromCookie(c)
	if err != nil || user.ID == 0 {
		c.Redirect(http.StatusFound, loginPageLocation)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Validate Campaign
	var campaign models.Campaign
	id := c.Params.ByName("id")
	models.DB.Preload("Users").Preload("Lead").First(&campaign, id)
	if campaign.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// See if user is linked to campaign or is Lead
	hasConnection := campaign.LeadID == user.ID
	for _, campaignUser := range campaign.Users {
		hasConnection = hasConnection || campaignUser.ID == user.ID
	}
	if !hasConnection {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.Set("user", user)
	c.Set("campaign", campaign)

	c.Next()
}

func RequireAuthAndCampaignAccess(c *gin.Context) {
	filePathParts := strings.Split(strings.TrimPrefix(c.Param("filepath"), "/"), "/")

	if len(filePathParts) < 3 {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// Validate User
	user, err := retrieveUserFromCookie(c)
	if err != nil || user.ID == 0 {
		c.Redirect(http.StatusFound, loginPageLocation)
		c.AbortWithStatus(http.StatusNotFound)
	}

	// Validate Campaign @todo;
	var campaign models.Campaign
	id := filePathParts[0]
	models.DB.Preload("Users").Preload("Lead").First(&campaign, id)
	if campaign.ID == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// See if user is linked to campaign or is Lead
	isLead := campaign.LeadID == user.ID
	hasConnection := false
	for _, campaignUser := range campaign.Users {
		hasConnection = hasConnection || campaignUser.ID == user.ID
	}

	// Only allow access to images if not a lead
	if !isLead && filePathParts[1] != "images" {
		c.AbortWithStatus(http.StatusNotFound)
	} else if !isLead && !hasConnection {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.Set("user", user)
	c.Set("campaign", campaign)

	c.Next()
}

func OptionalAuth(c *gin.Context) {
	user, err := retrieveUserFromCookie(c)
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
				c.Redirect(302, loginPageLocation)
				c.AbortWithStatus(401)
			}

			// Get user, if exists
			models.DB.First(&user, claims["ID"])
		}
	}
	return user, err
}
