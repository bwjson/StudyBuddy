package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
)

// @Summary      Get tags by UserID
// @Description  Get user's tag information by user ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
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

// @Summary      Get sort Order
// @Description  Retrieve a sort list of all orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   successResponse
// @Failure      500  {object}  errorResponse
// @Router       /user/sort_by&sort_order [get]

func (h *Handler) getSortOrder(c *gin.Context) (string, error) {
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	
	
	if sortBy != "id" && sortBy != "name" && sortBy != "title" && sortBy != "description" && sortBy != "username"  {
		return "", fmt.Errorf("invalid sort_by parameter")
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		return "", fmt.Errorf("invalid sort_order parameter")
	}
	return sortBy + " " + sortOrder, nil
}
// @Summary      Get pagination Order
// @Description  Retrieve a pagination list of all orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200  {array}   successResponse
// @Failure      500  {object}  errorResponse
// @Router       /user/{id}?page [get]
func (h *Handler) getPagination(c *gin.Context) (int, int, error) {
	page := c.DefaultQuery("page", "1")
	limit := 10
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1  {
		return 0, 0, fmt.Errorf("invalid page number")
	}
	offset := (pageNum - 1) * limit
	return offset, limit, nil
}

// @Summary      Get User by tagID
// @Description  Get tags information by user ID
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Success      200  {object}  successResponse
// @Failure      404  {object}  errorResponse
// @Router       /user/tag/{id} [get]
func (h *Handler) getUsersByTag(c *gin.Context) {
	tagID := c.Param("id")
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

	if len(users) == 0 {
		NewErrorResponse(c, http.StatusNotFound, "No users found for the specified tag")
		return
	}

	h.log.Info("getUsersByTag handler: Users successfully retrieved")

	NewSuccessResponse(c, http.StatusOK, "Users successfully retrieved", users)
}