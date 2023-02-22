package handler

import (
	"api/pkg/entities"
	"encoding/json"
	"github.com/Shopify/sarama"
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

	// Create a new Kafka producer
	producer, err := sarama.NewAsyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		logrus.Errorf("Error creating Kafka producer: %s", err.Error())
		return
	}

	// Convert the input to a JSON string
	data, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("Error marshaling input to JSON: %s", err.Error())
		return
	}

	// Create a new Kafka message containing the JSON data
	msg := &sarama.ProducerMessage{
		Topic: "card-creation",
		Value: sarama.StringEncoder(data),
	}

	// Send the message to Kafka
	producer.Input() <- msg

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
