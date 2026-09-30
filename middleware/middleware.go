package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"website-api/library/response"
)

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role_name") != "super_admin" {
			// Pakai library/response supaya bentuk error konsisten dengan handler
			// yang lain; AbortWithStatusJSON menghasilkan bentuk JSON yang beda.
			response.Error(c, http.StatusForbidden, "only super admin can access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
