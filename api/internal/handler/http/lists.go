package http

import (
	"api/pkg/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Create cards List
// @Security ApiKeyAuth
// @Tags lists
// @Description create cards Lists
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body entities.Lists true "list info"
// @Success 200 {integer} integer 1
// @Router /api/list [post]
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

// @Summary Create cards List
// @Security ApiKeyAuth
// @Tags lists
// @Description create cards ListData
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body ListData true "list info"
// @Success 200 {integer} integer 1
// @Router /api/list [get]
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

// @Summary Get a list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description create cards Response
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body entities.Lists true "list info"
// @Success 200 {integer} integer 1
// @Router /api/list/{id}  [get]
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

// @Summary Update a list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description create cards Response
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body entities.Lists true "list info"
// @Success 200 {integer} integer 1
// @Router /api/list/{id}  [post]
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

// @Summary Update a list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description create cards Response
// @ID create-list
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Router /api/list/{id}  [delete]
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
