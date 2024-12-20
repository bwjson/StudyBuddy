package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getUserTags(c *gin.Context) {

}

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

	h.log.Info("getAllTags handler: Successfully retrieved all tafs")

	NewSuccessResponse(c, http.StatusOK, "Successfully retrieved all tags", tags)
}

func (h *Handler) getUsersByTag(c *gin.Context) {
	tagID := c.Param("id")
	var tag dto.Tag
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_Order", "asc")
	page := c.DefaultQuery("page", "1")
	limit := 5

	if sortOrder != "asc" && sortOrder != "desc" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid sort_order parameter")
		return
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid page number")
		return
	}
	offset := (pageNum - 1) * limit
	order := sortBy + " " + sortOrder

	if err := h.db.First(&tag, tagID).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "Tag not found")
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

	if len(users) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No users found for the specified tag")
		return
	}

	h.log.Info("getUsersByTag handler: Users successfully retrieved")

	NewSuccessResponse(c, http.StatusOK, "Users successfully retrieved", users)
	
}