package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, statusCode int, message string, data any) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	ctx.JSON(statusCode, response{
		Data:    data,
		Message: message,
	})
}
