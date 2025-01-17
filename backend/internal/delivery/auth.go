package delivery

import (
	"fmt"
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/bwjson/StudyBuddy/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary      Sign in
// @Description  Sign in
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.SignInInput  true  "Sign up data"
// @Success      200   {object}  successResponse
// @Failure      400   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input dto.SignInInput
	var user dto.User
	var response dto.TokenResponse

	if err := c.BindJSON(&input); err != nil {
		h.log.Error("signIn handler: could not parse data from user")
		NewErrorResponse(c, http.StatusBadRequest, "signIn handler: could not parse data from user")
		return
	}

	input.Password = pkg.GenerateHashedPassword(input.Password)

	result := h.db.First(&user, "email = ? AND password_hash = ?", input.Email, input.Password)

	if result.Error != nil {
		h.log.Error("signIn handler: could not find user")
		NewErrorResponse(c, http.StatusInternalServerError, "signIn handler: could not find user")
		return
	}

	// AccessToken generating
	accessToken, err := pkg.TokenGen(input.Email)
	if err != nil {
		h.log.Error("signIn handler: could not generate access token")
		NewErrorResponse(c, http.StatusInternalServerError, "signIn handler: could not generate access token")
		return
	}

	// RefreshToken generating
	refreshToken, err := pkg.RefreshTokenGen()
	if err != nil {
		h.log.Error("signIn handler: could not generate refresh token")
		NewErrorResponse(c, http.StatusInternalServerError, "signIn handler: could not generate refresh token")
		return
	}

	// Response
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	NewSuccessResponse(c, http.StatusOK, "Successfully authenticated", response)
}

// @Summary      Sign up
// @Description  Sign up using email verification
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      dto.User  true  "Sign up data"
// @Success      200   {object}  successResponse
// @Failure      400   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input dto.User

	if err := c.BindJSON(&input); err != nil {
		h.log.Info("signUp handler: Failed to bind user input")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// hashing the password
	input.PasswordHash = pkg.GenerateHashedPassword(input.PasswordHash)

	// generating token
	token, err := pkg.TokenGen(input.Email)

	if err != nil {
		h.log.Info("signUp handler: Failed to generate token")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	input.VerificationToken = token

	if err := h.db.Create(&input).Error; err != nil {
		h.log.Info("signUp handler: Failed to create user")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	url := fmt.Sprintf("https://studybuddy-l0c9.onrender.com/auth/%s", token)

	h.smtp.SendVerifyingEmail(input.Email, "Registration", url)

	NewSuccessResponse(c, http.StatusOK, "Please, check your email address and click verifying link", nil)
}

func (h *Handler) verifyEmail(c *gin.Context) {
	var user dto.User
	token := c.Param("token")

	if err := h.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		h.log.Info("verifyEmail handler: User not found")
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.IsActive = true
	user.VerificationToken = ""

	if err := h.db.Save(&user).Error; err != nil {
		h.log.Info("verifyEmail handler: Failed to save user")
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	NewSuccessResponse(c, http.StatusOK, "Successfully verified user", nil)
}
