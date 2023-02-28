package main

import (
	"encoding/json"
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

	partitionConsumer := handlers.InitKafka()

	// Consume messages from the Kafka topic
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var input Card
			err := json.Unmarshal([]byte(string(msg.Value)), &input)
			if err != nil {
				logrus.Errorf("Error parsing kafka data: %s", err)
			}
			err = handlers.SetImageToCard(input.Id, input.Front)
			if err != nil {
				logrus.Errorf("Error generating image: %s", err.Error())
			}

			err = handlers.SetTranslateToCard(input.Id, input.Front)
			if err != nil {
				logrus.Errorf("Error generating translate: %s", err.Error())
			}
			break
		case err := <-partitionConsumer.Errors():
			logrus.Errorf("Error recieve data from kafka queue: %s", err)
			break

		}
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
