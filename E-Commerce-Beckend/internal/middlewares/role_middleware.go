package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredusertype string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usertype, exists := c.Get("userType")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "user type not found"})
			c.Abort()
			return
		}
		if usertype != requiredusertype {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}
