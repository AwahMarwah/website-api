package middleware

import "github.com/gin-gonic/gin"

func NgrokSkipWarning() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Header.Get("ngrok-skip-browser-warning") == "" {
			c.Request.Header.Set("ngrok-skip-browser-warning", "true")
		}

		c.Next()
	}
}
