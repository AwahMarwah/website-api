package review

import (
	"net/http"
	"website-api/library/pagination"
	"website-api/library/response"
	reviewModel "website-api/model/review"

	"github.com/gin-gonic/gin"
)

// @Summary List Reviews
// @Description Ambil daftar ulasan produk (public)
// @Tags Review
// @Produce json
// @Param id path string true "Product ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} review.SwaggerReviewList "Berhasil mengambil ulasan"
// @Router /product/{id}/reviews [get]
func (c *controller) ListByProduct(ctx *gin.Context) {
	var reqQuery reviewModel.ListReviewReqQuery
	if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqQuery.Offset = pagination.Offset(&reqQuery.Limit, &reqQuery.Page)

	productID := ctx.Param("id")
	resData, count, statusCode, err := c.reviewService.ListByProduct(productID, reqQuery.Page, reqQuery.Limit, reqQuery.Offset)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.SuccessWithPage(ctx, statusCode, "", resData, reqQuery.Page, reqQuery.Limit, count)
}