package KafkaRPC

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"imageGeneration/internal/service"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitKafka() sarama.PartitionConsumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatal("Error creating consumer: ", err)
	}

	partitionConsumer, err := consumer.ConsumePartition("card-creation", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error creating partition consumer: ", err)
	}

	return partitionConsumer
}

func (h *Handler) SetImageToCard(cardID int, prompt string) error {
	req, err := createRequest(prompt)
	if err != nil {
		return err
	}

	headers := createHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	err = h.service.SetImageToCard(cardID, response.Data[0].URL)
	if err != nil {
		return err
	}
	return nil

}
