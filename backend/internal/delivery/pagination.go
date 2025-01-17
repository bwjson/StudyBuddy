package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (h *Handler) getSortOrder(c *gin.Context) (string, error) {
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	if sortBy != "id" && sortBy != "name" && sortBy != "title" && sortBy != "description" && sortBy != "username" {
		return "", fmt.Errorf("invalid sort_by parameter")
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		return "", fmt.Errorf("invalid sort_order parameter")
	}
	return sortBy + " " + sortOrder, nil
}

func (h *Handler) getPagination(c *gin.Context) (int, int, error) {
	page := c.DefaultQuery("page", "1")
	limit := 10
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		return 0, 0, fmt.Errorf("invalid page number")
	}
	offset := (pageNum - 1) * limit
	return offset, limit, nil
}
