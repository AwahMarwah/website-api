package menu

import (
	"website-api/library/response"

	"github.com/gin-gonic/gin"
)

func (c *controller) Tree(ctx *gin.Context) {
	resData, statusCode, err := c.menuService.Tree()
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}
