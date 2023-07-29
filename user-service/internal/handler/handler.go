package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stas-bukovskiy/wish-scribe/packages/logger"
	_ "github.com/stas-bukovskiy/wish-scribe/user-service/docs"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	logger   logger.Logger
}

func NewHandler(services *service.Service, logger logger.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api/v1")
	{
		users := api.Group("/users", h.userIndemnity)
		{
			users.GET("/:id", h.getById)
		}
		tokens := api.Group("/tokens")
		{
			tokens.GET("/:token", h.verifyToken)
		}
	}
	return router
}
