package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/stas-bukovskiy/wish-scribe/user-service/docs"
	"github.com/stas-bukovskiy/wish-scribe/user-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.singUp)
		auth.POST("/sign-in", h.singIn)
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
