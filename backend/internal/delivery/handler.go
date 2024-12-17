package delivery

import (
	_ "github.com/bwjson/StudyBuddy/docs"
	"github.com/bwjson/StudyBuddy/pkg/smtp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	db   *gorm.DB
	log  *logrus.Logger
	smtp *smtp.SMTPServer
}

func NewHandler(db *gorm.DB, log *logrus.Logger, smtp *smtp.SMTPServer) *Handler {
	return &Handler{db: db, log: log, smtp: smtp}
}

// @title           StudyBuddy API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://example.com"}, // Укажите разрешённые источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rate limit middleware (request per second = 1, max = 3)
	router.Use(RateLimiter(1, 3))

	user := router.Group("/user")
	{
		user.POST("/", h.createUser)
		user.GET("/", h.getAllUsers)
		user.GET("/:id", h.getUserByID)
		user.PUT("/:id", h.updateUserByID)
		user.DELETE("/:id", h.deleteUserByID)
		user.POST("/email", h.sendEmail)
	}

	return router
}
