package order

import (
	"encoding/json"
	"io"
	"net/http"
	"website-api/library/response"
	midtransProvider "website-api/third-party/provider/midtrans"

	"github.com/gin-gonic/gin"
)

// @Summary Handle Payment Notification
// @Description Merespons notifikasi webhook dari Midtrans (public endpoint, tanpa auth)
// @Tags 4. Order
// @Accept json
// @Produce json
// @Param req body midtrans.NotificationPayload true "Midtrans Notification Payload"
// @Success 200 {object} order.NotificationResponse "Notifikasi diproses"
// @Router /order/notification [post]
func (c *controller) HandleNotification(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "failed to read request body")
		return
	}

	var payload midtransProvider.NotificationPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid notification payload")
		return
	}

	// verifikasi signature via provider
	if !c.midtrans.VerifySignature(payload) {
		response.Error(ctx, http.StatusBadRequest, "invalid signature")
		return
	}

	resData, statusCode, err := c.orderService.HandleNotification(payload)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}

	// Midtrans mengharapkan 200 dengan body "ok" agar tidak retry
	response.Success(ctx, statusCode, "", resData)
}
