package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessWithPage(ctx *gin.Context, statusCode int, message string, data any, current, size int, total int64) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	ctx.JSON(statusCode, responseWithPage{
		Data:    data,
		Message: message,
		Page: page{
			Current: current,
			Size:    size,
			Total:   total,
		},
	})
}
