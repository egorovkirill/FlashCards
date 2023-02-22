package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

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
			fmt.Println("Received message: ", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Println("Error: ", err)
			return
		}
	}
}
