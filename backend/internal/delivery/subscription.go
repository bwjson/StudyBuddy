package delivery

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Security ApiKeyAuth
// @Summary      Buy subscription
// @Description  Subscription by PayPal
// @Tags         paypal
// @Accept       json
// @Produce      json
// @Param        request body buySubscriptionRequest true "Subscription Data"
// @Success      200  {object}  successGRPCResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Router       /subscriptions/buy/ [post]
func (h *Handler) buySubscription(c *gin.Context) {
	var req buySubscriptionRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Invalid request body"})
		return
	}

	ctx := context.Background()
	resp, err := h.client.BuySubscription(ctx, req.Email, req.CardNumber, req.ValidUntil, req.Cvv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, successGRPCResponse{Message: resp.Response, Detail: resp.Detail})

}

// @Security ApiKeyAuth
// @Summary      Cancel subscription
// @Description  Cancel subscription via PayPal
// @Tags         paypal
// @Accept       json
// @Produce      json
// @Param        request body cancelSubscriptionRequest true "Subscription Cancellation Data"
// @Success      200  {object}  successGRPCResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Router       /subscriptions/cancel/ [post]
func (h *Handler) cancelSubscription(c *gin.Context) {
	var req cancelSubscriptionRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Invalid request body"})
		return
	}

	ctx := context.Background()
	resp, err := h.client.CancelSubscription(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, successGRPCResponse{Message: resp.Response, Detail: resp.Detail})
}

// buySubscriptionRequest структура для приема JSON
type buySubscriptionRequest struct {
	CardNumber string `json:"card_number" binding:"required"`
	Cvv        string `json:"cvv" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	ValidUntil string `json:"valid_until" binding:"required"`
}

// successResponse структура ответа для успешного запроса
type successGRPCResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type cancelSubscriptionRequest struct {
	CardNumber string `json:"card_number" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}
