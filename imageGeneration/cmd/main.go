package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"io"
	"log"
	"net/http"
)

type Card struct {
	Front string `json:"front"`
}

type Response struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type Request struct {
	Prompt         string `json:"prompt"`
	NumImages      int    `json:"num_images"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

var reqHeaders = map[string]string{
	"Content-Type":  "application/json",
	"Authorization": "Bearer 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",
}

var apiURL = "https://api.openai.com/v1/images/generations"

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Replace "localhost:9092" with your Kafka broker's address
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatal("Error creating consumer: ", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal("Error closing consumer: ", err)
		}
	}()

	// Replace "test" with the name of your Kafka topic
	partitionConsumer, err := consumer.ConsumePartition("card-creation", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error creating partition consumer: ", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal("Error closing partition consumer: ", err)
		}
	}()

	// Consume messages from the Kafka topic
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// Replace this with your actual processing logic
			var input Card
			fmt.Println("Received message: ", string(msg.Value))
			err := json.Unmarshal([]byte(string(msg.Value)), &input)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(GenerateImage(input.Front))
		case err := <-partitionConsumer.Errors():
			fmt.Println("222222o")
			log.Println("Error: ", err)
			return
		}
	}
}

func GenerateImage(prompt string) (string, error) {
	request := Request{
		Prompt:         prompt,
		NumImages:      1,
		Size:           "256x256",
		ResponseFormat: "url",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	for key, value := range reqHeaders {
		req.Header.Set(key, value)
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Data[0].URL, nil
}
