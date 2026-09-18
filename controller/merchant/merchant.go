package merchant

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	merchantModel "website-api/model/merchant"
	merchantRepo "website-api/repository/merchant"
	merchantService "website-api/service/merchant"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type controller struct {
	merchantService merchantService.IService
}

func NewController(db *gorm.DB) *controller {
	return &controller{merchantService: merchantService.NewService(merchantRepo.NewRepo(db))}
}

func (c *controller) Register(ctx *gin.Context) {
	var req merchantModel.CreateMerchantReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	statusCode, err := c.merchantService.Register(&req, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

func (c *controller) GetMy(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	resData, statusCode, err := c.merchantService.GetMy(userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", resData)
}

func (c *controller) UpdateMy(ctx *gin.Context) {
	var req merchantModel.UpdateMerchantReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	statusCode, err := c.merchantService.UpdateMy(userInfo.UserID, &req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}

func (c *controller) Approve(ctx *gin.Context) {
	var reqPath merchantModel.ReqPath
	if err := ctx.ShouldBindUri(&reqPath); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.merchantService.Approve(reqPath.Id)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}