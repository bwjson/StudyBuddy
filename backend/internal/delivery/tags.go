package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary      Get tags by UserID
// @Description  Get user's tag information by user ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        sort_by     query   string  false  "Sort by field (id, title, description)"
// @Param        sort_order  query   string  false  "Sort order (asc, desc)"
// @Param        page  query     int     false "Page number"
// @Success      200  {object}  successResponse
// @Failure      404  {object}  errorResponse
// @Router       /user/usertags/{id} [get]
func (h *Handler) getTagsByUser(c *gin.Context) {
	userID := c.Param("id")
	var tags []dto.Tag
	var user []dto.User
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

	if err := h.db.Joins("JOIN user_tags ut ON ut.tag_id = tags.id").
		Where("ut.user_id = ?", userID).
		Order(order).
		Limit(limit).
		Offset(offset).
		Find(&tags).Error; err != nil {
		h.log.Error("getTagsByUser handler: Failed to get user's tags, userID:", userID, "error:", err)
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user's tags")
		return
	}

	if err := h.db.First(&user, userID).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	if len(tags) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No tags found for the specified user")
		return
	}

	h.log.Info("getTagsByUser handler: Successfully retrieved tags for userID:", userID)

	NewSuccessResponse(c, http.StatusOK, "Successfully retrieved user's tags", tags)
}

// @Summary      Get all tags
// @Description  Retrieve a list of all tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {array}   successResponse
// @Failure      500  {object}  errorResponse
// @Router       /user/tags [get]
func (h *Handler) getAllTags(c *gin.Context) {
	var tags []dto.Tag

	if err := h.db.Find(&tags).Error; err != nil {
		h.log.Error("getAllTags handler: Failed to get all tags")
		NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tags")
		return
	}

	if len(tags) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No tags found")
		return
	}

	h.log.Info("getAllTags handler: Successfully retrieved all tags")

	NewSuccessResponse(c, http.StatusOK, "Successfully retrieved all tags", tags)
}

// @Summary      Get User by tagID
// @Description  Get tags information by user ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Param        sort_by     query   string  false  "Sort by field (id, name, username)"
// @Param        sort_order  query   string  false  "Sort order (asc, desc)"
// @Param        page  query     int     false "Page number"
// @Success      200  {object}  successResponse
// @Failure      404  {object}  errorResponse
// @Router       /user/tag/{id} [get]
func (h *Handler) getUsersByTag(c *gin.Context) {
	tagID := c.Param("id")
	var totalUsers []dto.User
	var tag dto.Tag

	if err := h.db.First(&tag, tagID).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "Tag not found")
		return
	}

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

	var users []dto.User
	if err := h.db.Joins("JOIN user_tags ut ON ut.user_id = users.id").
		Where("ut.tag_id = ?", tag.ID).
		Order(order).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "Users not found")
		return
	}

	if err := h.db.Joins("JOIN user_tags ut ON ut.user_id = users.id").
		Where("ut.tag_id = ?", tag.ID).
		Find(&totalUsers).Error; err != nil {
		return
	}

	totalCount := len(totalUsers)

	if len(users) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No users found for the specified tag")
		return
	}

	response := dto.UsersWithPagination{
		User:       users,
		TotalCount: totalCount,
	}

	h.log.Info("getUsersByTag handler: Users successfully retrieved")

	NewSuccessResponse(c, http.StatusOK, "Users successfully retrieved", response)
}
