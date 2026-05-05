package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role_name")

		fmt.Println(role, "ini role nya")

		if role != "super_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "only super admin can access this resource",
			})
			return
		}

		c.Next()
	}
}
