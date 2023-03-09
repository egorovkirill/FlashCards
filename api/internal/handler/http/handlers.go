package http

import (
	"api/internal/service"
	_ "github.com/egorovkirill/FlashCards/tree/master/api/cmd/docs"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	route := gin.New()
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.Use(PrometheusMiddleware())
	metrics := route.Group("/prometheus")
	{
		metrics.GET("/metrics", func(c *gin.Context) {
			promhttp.Handler().ServeHTTP(c.Writer, c.Request)
		})
	}
	auth := route.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUP)
		auth.POST("/sign-in", h.SignIN)
	}

	api := route.Group("/api", h.userIdentity)
	{
		lists := api.Group("/list")
		{
			lists.POST("/", h.CreateList)
			lists.GET("/", h.GetLists)
			lists.GET("/:id", h.GetListById)
			lists.POST("/:id", h.UpdateListById)
			lists.DELETE("/:id", h.DeleteListById)
		}
		items := api.Group("/:listID/card")
		{
			items.POST("/", h.CreateCard)
			items.GET("/", h.GetCardsInList)
			items.GET("/:cardID", h.GetCardById)
			items.DELETE("/:cardID", h.DeleteCardById)
		}

	}
	return route
}
