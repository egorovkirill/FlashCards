package repository

import (
	"ToDo/pkg/entities"
	"github.com/jmoiron/sqlx"
)

type Authenthication interface {
	CreateUser(user entities.User) (int, error)
	ValidateUser(user entities.User) (entities.User, error)
}

type WordsList interface {
	CreateList(userId int, title string) (int, error)
	GetLists(userId int) ([]entities.Lists, error)
	GetListById(userId, ListId int) ([]entities.Lists, error)
	UpdateListById(userId, ListId int, title string) error
	DeleteListById(userId, ListId int) error
}

type ListCards interface {
	CreateCard(listID int, cards entities.Cards) (int, error)
	GetCardsInList(userId, listID int) ([]entities.Cards, error)
	GetCardById(userID, listID, itemID int) ([]entities.Cards, error)
	DeleteCardById(userID, itemID int) error
}

type Repository struct {
	Authenthication
	WordsList
	ListCards
}

func NewPostgresRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authenthication: NewAuthPostgres(db),
		WordsList:       NewCreateListPostgres(db),
		ListCards:       NewCardPostgres(db),
	}
}
