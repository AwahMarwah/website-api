package review

import (
	"net/http"
	"website-api/library/response"
	"website-api/middleware"
	reviewModel "website-api/model/review"

	"github.com/gin-gonic/gin"
)

// @Summary Create Review
// @Description Tambah ulasan produk (hanya user dengan order COMPLETED utk produk tsb)
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param req body review.CreateReviewReq true "Review Request"
// @Success 201 {object} map[string]interface{} "Review berhasil dibuat"
// @Router /product/{id}/reviews [post]
func (c *controller) Create(ctx *gin.Context) {
	userInfo, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	var req reviewModel.CreateReviewReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	req.ProductID = ctx.Param("id")

	statusCode, err := c.reviewService.Create(&req, userInfo.UserID)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Success(ctx, statusCode, "", nil)
}