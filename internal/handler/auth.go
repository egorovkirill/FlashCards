package handler

import (
	"ToDo/pkg/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) SignUP(c *gin.Context) {
	var input entities.User
	if err := c.BindJSON(&input); err != nil {
		return
	}
	id, err := h.service.CreateUser(input)
	if err != nil {
		logrus.Error(err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) SignIN(c *gin.Context) {
	var input entities.User
	if err := c.BindJSON(&input); err != nil {
		logrus.Error("Error while parsing user credentials", err.Error())
	}

	token, err := h.service.GenerateToken(input)
	if err != nil {
		logrus.Error(err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
