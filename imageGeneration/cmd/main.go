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
	"sync"
	"time"
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

	var wg sync.WaitGroup
	numWorkers := 10 // Number of worker goroutines
	jobCh := make(chan []byte)
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for job := range jobCh {
				var input Card
				err := json.Unmarshal(job, &input)
				if err != nil {
					logrus.Errorf("Error parsing kafka data: %s", err)
					continue
				}
				err = handlers.SetImageToCard(input.Id, input.Front)
				if err != nil {
					logrus.Errorf("Error generating image: %s", err.Error())
					continue
				}
				err = handlers.SetTranslateToCard(input.Id, input.Front)
				if err != nil {
					logrus.Errorf("Error generating translate: %s", err.Error())
					continue
				}
			}
			wg.Done()
		}()
	}

	// Read messages from the Kafka topic and send them to the worker pool
	counter := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			jobCh <- msg.Value
			counter++
			if counter == 49 {
				time.Sleep(80 * time.Second)
				counter = 0
			}
		case err := <-partitionConsumer.Errors():
			logrus.Errorf("Error receiving data from kafka queue: %s", err)
		}
	}
	close(jobCh)
	wg.Wait()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
