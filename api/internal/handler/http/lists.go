package http

import (
	"api/pkg/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) CreateList(c *gin.Context) {
	var input entities.Lists
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}

	if err := c.BindJSON(&input); err != nil {
		logrus.Error("Error parsing list data from user", err.Error())
		return
	}
	id, err := h.service.WordsList.CreateList(userId.(int), input.Title)
	if err != nil {
		logrus.Errorf("Error creating list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type ListData struct {
	Lists []entities.Lists `json:"lists"`
}

func (h *Handler) GetLists(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}
	response, err := h.service.GetLists(userId.(int))
	if err != nil {
		logrus.Errorf("Error getting lists: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, ListData{
		Lists: response,
	})

}

type Response struct {
	Response []entities.Lists `json:"response"`
}

func (h *Handler) GetListById(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("id")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing param")
		return
	}
	response, err := h.service.GetListById(userId.(int), listId)
	c.JSON(http.StatusOK, Response{
		Response: response,
	})
}

func (h *Handler) UpdateListById(c *gin.Context) {
	var input entities.Lists
	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("Error parsing input: %s", err.Error())
		return
	}
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("id")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing param")
		return
	}
	err = h.service.UpdateListById(userId.(int), listId, input.Title)
	if err != nil {
		logrus.Errorf("Error updating list: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) DeleteListById(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("id")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing param")
		return
	}
	err = h.service.DeleteListById(userId.(int), listId)
	if err != nil {
		logrus.Errorf("Error deleting list: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
