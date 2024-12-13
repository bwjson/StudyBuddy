package delivery

import (
	"github.com/bwjson/StudyBuddy/internal/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"github.com/gin-contrib/cors"
)

// TODO: Реализовать структуру хэндлера и конструктор с пробросами по слоям
type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
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

	// TODO: Сделать нормальное возвращение ошибок
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
	// TODO: Тут нужно использовать модельки
	var input dto.User

	if err := c.ShouldBindJSON(&input); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
func (h *Handler) getAllUsers(c *gin.Context) {
    var users []dto.User 

    if err := h.db.Find(&users).Error; err != nil {
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

    c.JSON(http.StatusOK, gin.H{
        "users": users,
    })
}


func (h *Handler) updateUserByID(c *gin.Context) {
	id := c.Param("id")
	var input dto.User
	var user dto.User

	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON message",
		})
		return
	}

	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	user.Name = input.Name
	user.Username = input.Username
	user.PasswordHash = input.PasswordHash

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully updated",
		"user":    user,
	})
}

func (h *Handler) deleteUserByID(c *gin.Context) {
	id := c.Param("id")
	var user dto.User

	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
