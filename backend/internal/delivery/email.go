package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary      Send an email
// @Description  Send an email to the specified address
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        email  body      dto.EmailInput  true  "Email input"
// @Success      200    {object}  successResponse
// @Failure      400    {object}  errorResponse
// @Failure      500    {object}  errorResponse
// @Router       /user/email [post]
func (h *Handler) sendEmail(c *gin.Context) {
	var input dto.EmailInput

	if err := c.BindJSON(&input); err != nil {
		h.log.Error("sendEmail handler: Invalid JSON sent")
		NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON message")
		return
	}

	err := h.smtp.SendEmail(input.Email, input.Subject, input.Message)
	if err != nil {
		h.log.Error("Failed to sent an email")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to send email")
		return
	}

	h.log.Info("Email sent successfully")
	NewSuccessResponse(c, http.StatusOK, "Successfully sent email", nil)
}
