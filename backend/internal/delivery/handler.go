package delivery

import (
	_ "github.com/bwjson/StudyBuddy/docs"
	"github.com/bwjson/StudyBuddy/internal/grpc"
	"github.com/bwjson/StudyBuddy/pkg/smtp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	log    *logrus.Logger
	smtp   *smtp.SMTPServer
	client *grpc.Client
}

func NewHandler(db *gorm.DB, log *logrus.Logger, smtp *smtp.SMTPServer, client *grpc.Client) *Handler {
	return &Handler{db: db, log: log, smtp: smtp, client: client}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(RateLimiter(1, 3))

	user := router.Group("/user", h.adminIdentity)
	{
		user.POST("/", h.createUser)
		user.GET("/", h.getAllUsers)
		user.GET("/:id", h.getUserByID)
		user.PUT("/:id", h.updateUserByID)
		user.DELETE("/:id", h.deleteUserByID)
		user.POST("/email", h.sendEmail)
		//user.GET("/tag/:id", h.getUsersByTag)
		//user.GET("/tags", h.getAllTags)
		//user.GET("/usertags/:id", h.getTagsByUser)
	}

	tags := router.Group("/tags", h.userIdentity)
	{
		tags.GET("/:id", h.getUsersByTag)
		tags.GET("/", h.getAllTags)
		tags.GET("/usertags/:id", h.getTagsByUser)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.GET("/:token", h.verifyEmail)
	}

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/buy", h.buySubscription)
		subscriptions.POST("/cancel", h.cancelSubscription)
	}

	return router
}
