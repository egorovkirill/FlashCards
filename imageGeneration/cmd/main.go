package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"imageGeneration/internal/database"
	"imageGeneration/internal/handler/KafkaRPC"
	service "imageGeneration/internal/service"
	"log"
	"os"
)

type Card struct {
	Id           int    `json:"id"`
	Front        string `json:"front"`
	Back         string `json:"back"`
	ImageLink    string `json:"imageLink"`
	VoiceMessage string `json:"voiceMessage"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := initConfig(); err != nil {
		log.Fatal("Error loading config")
	}

	db := database.ConnectToPostgresDB(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	repository := database.NewRepository(db)
	services := service.NewService(repository)
	handlers := KafkaRPC.NewHandler(services)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatal("Error creating consumer: ", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal("Error closing consumer: ", err)
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition("card-creation", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error creating partition consumer: ", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal("Error closing consumer: ", err)
		}
	}()

	// Consume messages from the Kafka topic
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var input Card
			err := json.Unmarshal([]byte(string(msg.Value)), &input)
			if err != nil {
				logrus.Errorf("Error parsing kafka data: %s", err)
				return
			}
			err = handlers.SetImageToCard(input.Id, input.Front)
			if err != nil {
				logrus.Errorf("Error generating image: %s", err.Error())
				return
			}

		case err := <-partitionConsumer.Errors():
			logrus.Errorf("Error recieve data from kafka queue: %s", err)
			continue
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
