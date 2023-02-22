package handler

import (
	"api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	route := gin.New()
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
