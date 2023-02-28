package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectToPostgresDB(cfg Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode))

	if err != nil {
		logrus.Fatalf("Error initializing DB: %s", err.Error())
	}

	return db
}
