package response

import "github.com/gin-gonic/gin"

func Error(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, response{
		Data:    nil,
		Message: message,
	})
}
