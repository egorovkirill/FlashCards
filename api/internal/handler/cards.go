package handler

import (
	"ToDo/pkg/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ResponseCards struct {
	Response []entities.Cards `json:"response"`
}

func (h *Handler) CreateCard(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("listID")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing Param from URL: %s", err.Error())
		return
	}

	var input entities.Cards

	if err := c.BindJSON(&input); err != nil {
		logrus.Errorf("Error parsing input: %s", err.Error())
		return
	}

	_, err = h.service.CreateCard(listId, input)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, http.StatusOK)
}

func (h *Handler) GetCardsInList(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("listID")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing Param from URL: %s", err.Error())
		return
	}
	response, err := h.service.GetCardsInList(userId.(int), listId)
	if err != nil {
		logrus.Errorf("Error parsing cards in list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, ResponseCards{
		Response: response,
	})

}

func (h *Handler) GetCardById(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("listID")
	listId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing Param from URL: %s", err.Error())
		return
	}
	id = c.Param("cardID")
	cardId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing Param from URL: %s", err.Error())
		return
	}

	response, err := h.service.GetCardById(userID.(int), listId, cardId)
	if err != nil {
		logrus.Errorf("Error parsing cards in list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, ResponseCards{
		Response: response,
	})

}

func (h *Handler) DeleteCardById(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		return
	}
	id := c.Param("cardID")
	cardId, err := strconv.Atoi(id)
	if err != nil {
		logrus.Errorf("Error parsing Param from URL: %s", err.Error())
		return
	}

	err = h.service.DeleteCardById(userID.(int), cardId)
	if err != nil {
		logrus.Errorf("Error parsing cards in list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, http.StatusOK)
}
