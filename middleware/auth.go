package middleware

import (
	"net/http"

	"recruitment-portal/models"

	"github.com/gin-gonic/gin"
)

// AuthRequired ensures the user is logged in.
// It assumes that middleware (or the OAuth callback) has set a "user" key in Gin's context.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		// Optionally, you can cast sessionUser to *models.User here.
		c.Next()
	}
}

// SuperAdminRequired ensures that the logged in user is a Super Admin.
func SuperAdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		user, ok := sessionUser.(*models.User)
		if !ok || user.Role != models.RoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Super Admin privileges required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
