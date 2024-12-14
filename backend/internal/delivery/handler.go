package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

// TODO: Реализовать структуру хэндлера и конструктор с пробросами по слоям
type Handler struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewHandler(db *gorm.DB, log *logrus.Logger) *Handler {
	return &Handler{db: db, log: log}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://example.com"}, // Укажите разрешённые источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	example := router.Group("/example")
	{
		example.GET("/", h.exampleGreet)
		example.POST("/", h.exampleSend)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/create", h.createUser)
		auth.GET("/users", h.getAllUsers)
		auth.GET("/user/:id", h.getUserByID)
		auth.PUT("/user/:id", h.updateUserByID)
		auth.DELETE("/user/:id", h.deleteUserByID)
	}

	return router
}

func (h *Handler) exampleGreet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, from server! This is JSON data for Postman.",
	})
}

func (h *Handler) exampleSend(c *gin.Context) {
	var input dto.JsonInput

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON message",
		})
		return
	}

	if input.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The 'message' field is required and cannot be empty",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data successfully received",
	})
}

func (h *Handler) createUser(c *gin.Context) {
	var input dto.User

	if err := c.BindJSON(&input); err != nil {
		h.log.Error("createUser handler: Invalid JSON sent")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON message",
		})
		return
	}

	user := dto.User{
		Name:         input.Name,
		Username:     input.Username,
		PasswordHash: input.PasswordHash,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	h.log.Info("createUser handler: User successfully created")
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully created",
		"user":    user,
	})
}

func (h *Handler) getUserByID(c *gin.Context) {
	id := c.Param("id")
	var user dto.User

	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	h.log.Info("getUserByID handler: User successfully retrieved")

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	var users []dto.User

	if err := h.db.Find(&users).Error; err != nil {
		h.log.Error("getAllUsers handler: Failed to get all users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No users found",
		})
		return
	}

	h.log.Info("getAllUsers handler: Successfully retrieved all users")

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (h *Handler) updateUserByID(c *gin.Context) {
	id := c.Param("id")
	var input dto.User
	var user dto.User

	if err := c.Bind(&input); err != nil {
		h.log.Error("updateUserByID handler: Invalid JSON sent")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON message",
		})
		return
	}

	if err := h.db.First(&user, id).Error; err != nil {
		h.log.Error("updateUserByID handler: Failed to retrieve user")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user.Name = input.Name
	user.Username = input.Username
	user.PasswordHash = input.PasswordHash

	if err := h.db.Save(&user).Error; err != nil {
		h.log.Error("updateUserByID handler: Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	h.log.Info("updateUserByID handler: User successfully updated")

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully updated",
		"user":    user,
	})
}

func (h *Handler) deleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user dto.User

	if err := h.db.First(&user, id).Error; err != nil {
		h.log.Error("deleteUserByID handler: Failed to retrieve user")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		h.log.Error("deleteUserByID handler: Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	h.log.Info("deleteUserByID handler: User successfully deleted")

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
