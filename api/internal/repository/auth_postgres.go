package repository

import (
	"ToDo/pkg/entities"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (c *AuthPostgres) CreateUser(user entities.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id")
	row := c.db.QueryRow(query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *AuthPostgres) ValidateUser(user entities.User) (entities.User, error) {
	query := fmt.Sprintf("SELECT id, password, login FROM users WHERE login = $1")
	err := c.db.Get(&user, query, user.Login)
	if err != nil {
		return entities.User{}, err
	}
	return user, err
}
