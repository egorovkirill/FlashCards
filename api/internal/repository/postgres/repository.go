package postgres

import (
	"api/pkg/entities"
	"github.com/Shopify/sarama"
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
	CreateCard(listID int, cards entities.Cards) error
	GetCardsInList(userId, listID int) ([]entities.Cards, error)
	GetCardById(userID, listID, itemID int) ([]entities.Cards, error)
	DeleteCardById(userID, itemID int) error
	SetImageToCard(cardID int, image string) error
	SetTranslateToCard(cardID int, translate string) error
}

type Repository struct {
	Authenthication
	WordsList
	ListCards
}

func NewPostgresRepository(db *sqlx.DB, kafkaAsync sarama.AsyncProducer) *Repository {
	return &Repository{
		Authenthication: NewAuthPostgres(db),
		WordsList:       NewCreateListPostgres(db),
		ListCards:       NewCardPostgres(db, kafkaAsync),
	}
}
