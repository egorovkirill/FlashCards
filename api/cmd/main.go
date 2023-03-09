package main

import (
	handler "api/internal/handler/http"
	"api/internal/repository/postgres"
	"api/internal/service"
	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading local variables: %s", err.Error())
	}
	db := postgres.ConnectToPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	defer db.Close()

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	producer, err := sarama.NewAsyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		return
	}

	// Set the log retention to 7 days
	repos := postgres.NewPostgresRepository(db, producer)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	go handler.UpdateCPUMetrics()

	if err := router.Run(viper.GetString("PORT")); err != nil {
		logrus.Fatalf("Error occured while running server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
