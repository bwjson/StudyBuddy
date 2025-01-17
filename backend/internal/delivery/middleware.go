package delivery

import (
	"errors"
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/bwjson/StudyBuddy/pkg"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "email"
)

func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (h *Handler) userIdentity(c *gin.Context) {
	email, err := h.parseAuthHeader(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "Cannot parse auth header")
		h.log.Error(email)
		return
	}

	c.Set(userCtx, email)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	email, err := h.parseAuthHeader(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "Cannot parse auth header")
		h.log.Error(email)
		return
	}

	var user dto.User

	if err := h.db.First(&user, "email = ?", email).Error; err != nil {
		NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	if user.IsAdmin != true {
		NewErrorResponse(c, http.StatusUnauthorized, "User is not admin")
		return
	}

	c.Set(userCtx, email)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty auth header")
	}

	return pkg.Parse(headerParts[1])
}

func getEmailByContext(c *gin.Context, context string) (string, error) {
	emailFromCtx, ok := c.Get(context)
	if !ok {
		return "", errors.New("userCtx not found")
	}

	email, ok := emailFromCtx.(string)
	if !ok {
		return "", errors.New("userCtx is of invalid type")
	}

	return email, nil
}
