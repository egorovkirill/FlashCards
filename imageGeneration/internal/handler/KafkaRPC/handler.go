package KafkaRPC

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"imageGeneration/internal/service"
	"io"
	"log"
	"net/http"
	"os"
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

	partitionConsumer, err := consumer.ConsumePartition("cards_topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error creating partition consumer: ", err)
	}

	return partitionConsumer
}

func (h *Handler) SetImageToCard(cardID int, prompt string) error {
	req, err := GetImage(prompt)
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
	if len(response.Data) != 0 {
		err = h.service.SetImageToCard(cardID, response.Data[0].URL)
		if err != nil {
			return err
		}
	} else {
		return http.ErrBodyNotAllowed
	}

	return nil
}

func (h *Handler) SetTranslateToCard(cardID int, prompt string) error {
	url := fmt.Sprintf("https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text=%s", os.Getenv("YANDEX_API"), prompt)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var lookupResp LookupResponse
	err = json.NewDecoder(resp.Body).Decode(&lookupResp)
	if err != nil {
		return err
	}
	if len(lookupResp.Def) > 0 {
		if len(lookupResp.Def[0].Tr) > 0 {
			err = h.service.SetTranslateToCard(cardID, lookupResp.Def[0].Tr[0].Text)
			if err != nil {
				return err
			}
		}
	} else {
		return http.ErrBodyNotAllowed
	}

	return nil
}
