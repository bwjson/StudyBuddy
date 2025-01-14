package delivery

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetSortOrder(t *testing.T) {
	tests := []struct {
		name       string
		sortBy     string
		sortOrder  string
		expected   string
		expectErr  bool
	}{
		{"Valid sort parameters", "name", "asc", "name asc", false},
		{"Invalid sort_by parameter", "invalid", "asc", "", true},
		{"Invalid sort_order parameter", "id", "invalid", "", true},
		{"Valid sort parameters descending", "title", "desc", "title desc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request, _ = http.NewRequest("GET", "/test", nil)
			query := c.Request.URL.Query()
			query.Set("sort_by", tt.sortBy)
			query.Set("sort_order", tt.sortOrder)
			c.Request.URL.RawQuery = query.Encode()

			h := &Handler{}
			sortOrder, err := h.getSortOrder(c)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, sortOrder)
			}
		})
	}
}
