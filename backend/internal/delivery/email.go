package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
)

// @Summary      Send an email
// @Description  Send an email to the specified address
// @Tags         email
// @Accept       multipart/form-data
// @Produce      json
// @Param        email  formData  dto.EmailInput  true  "Email input"
// @Param        attachments  formData  file  false  "Attachments"
// @Success      200    {object}  successResponse
// @Failure      400    {object}  errorResponse
// @Failure      500    {object}  errorResponse
// @Router       /user/email [post]
func (h *Handler) sendEmail(c *gin.Context) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		NewErrorResponse(c, http.StatusInternalServerError, "Invalid .env file")
		return
	}

	var input dto.EmailInput

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		h.log.Error("sendEmail handler: Unable to parse form data")
		NewErrorResponse(c, http.StatusBadRequest, "Invalid form data")
		return
	}

	input.Email = c.PostForm("email")
	input.Subject = c.PostForm("subject")
	input.Message = c.PostForm("message")

	if input.Email == "" || input.Subject == "" || input.Message == "" {
		h.log.Error("sendEmail handler: Missing required fields")
		NewErrorResponse(c, http.StatusBadRequest, "Missing required fields")
		return
	}

	form := c.Request.MultipartForm
	files := form.File["attachments"]

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "khanapin65@gmail.com")
	mailer.SetHeader("To", input.Email)
	mailer.SetHeader("Subject", input.Subject)
	mailer.SetBody("text/plain", input.Message)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Failed to open file")
			return
		}
		defer file.Close()

		dst := "./uploads/" + fileHeader.Filename
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
			return
		}

		mailer.Attach(dst)
	}

	dialer := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_FROM"), os.Getenv("SMTP_PASSWORD"))
	err := dialer.DialAndSend(mailer)
	if err != nil {
		h.log.Error("Failed to send email")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to send email")
		return
	}

	h.log.Info("Email sent successfully")
	NewSuccessResponse(c, http.StatusOK, "Successfully sent email", nil)
}
