package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary      Create a new user
// @Description  Create a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.User  true  "User data"
// @Success      200   {object}  successResponse
// @Failure      400   {object}  errorResponse
// @Failure      500   {object}  errorResponse
// @Router       /user [post]
func (h *Handler) createUser(c *gin.Context) {
	var input dto.User

	if err := c.BindJSON(&input); err != nil {
		h.log.Error("createUser handler: Invalid JSON sent")
		NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON message")
		return
	}

	user := dto.User{
		Name:         input.Name,
		Username:     input.Username,
		PasswordHash: input.PasswordHash,
	}

	if err := h.db.Create(&user).Error; err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.log.Info("createUser handler: User successfully created")
	NewSuccessResponse(c, http.StatusCreated, "User successfully created", user)
}

// @Summary      Get user by ID
// @Description  Get user information by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  successResponse
// @Failure      404  {object}  errorResponse
// @Router       /user/{id} [get]
func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user dto.User

	if err := h.db.First(&user, id).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	h.log.Info("getUserByID handler: User successfully retrieved")

	NewSuccessResponse(c, http.StatusOK, "User successfully retrieved", user)
}

// @Summary      Get all users
// @Description  Retrieve a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   successResponse
// @Failure      500  {object}  errorResponse
// @Router       /user [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	var users []dto.User

	order, err := h.getSortOrder(c)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	
	offset, limit, err := h.getPagination(c)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.db.
		Order(order).
		Limit(limit).
		Offset(offset).Find(&users).Error; err != nil {
		h.log.Error("getAllUsers handler: Failed to get all users")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	if len(users) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No users found")
		return
	}

	h.log.Info("getAllUsers handler: Successfully retrieved all users")

	NewSuccessResponse(c, http.StatusOK, "Successfully retrieved all users", users)
}

// @Summary      Update user by ID
// @Description  Update user information by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int    true  "User ID"
// @Param        user body      dto.User  true  "Updated user data"
// @Success      200  {object}  successResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /user/{id} [put]
func (h *Handler) updateUserByID(c *gin.Context) {
	id := c.Param("id")
	var input dto.User
	var user dto.User

	if err := c.BindJSON(&input); err != nil {
		h.log.Error("updateUserByID handler: Invalid JSON sent")
		NewErrorResponse(c, http.StatusBadRequest, "Invalid JSON message")
		return
	}

	if err := h.db.First(&user, id).Error; err != nil {
		h.log.Error("updateUserByID handler: Failed to retrieve user")
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.Name = input.Name
	user.Username = input.Username
	user.PasswordHash = input.PasswordHash

	if err := h.db.Save(&user).Error; err != nil {
		h.log.Error("updateUserByID handler: Failed to update user")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.log.Info("updateUserByID handler: User successfully updated")

	NewSuccessResponse(c, http.StatusOK, "User successfully updated", user)
}

// @Summary      Delete user by ID
// @Description  Delete a user from the system by user ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  successResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /user/{id} [delete]
func (h *Handler) deleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user dto.User

	if err := h.db.First(&user, id).Error; err != nil {
		h.log.Error("deleteUserByID handler: Failed to retrieve user")
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		h.log.Error("deleteUserByID handler: Failed to delete user")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.log.Info("deleteUserByID handler: User successfully deleted")

	NewSuccessResponse(c, http.StatusOK, "User successfully deleted", nil)
}
