package service

import (
	"ToDo/internal/repository"
	"ToDo/pkg/entities"
)

type Authenthication interface {
	CreateUser(user entities.User) (int, error)
	GenerateToken(user entities.User) (string, error)
	ParseToken(accessToken string) (int, error)
}

type WordsList interface {
	CreateList(userId int, title string) (int, error)
	GetLists(userId int) ([]entities.Lists, error)
	GetListById(userId int, listId int) ([]entities.Lists, error)
	UpdateListById(userId, ListId int, title string) error
	DeleteListById(userId, ListId int) error
}

type ListCards interface {
	CreateCard(listID int, cards entities.Cards) (int, error)
	GetCardsInList(userId, listID int) ([]entities.Cards, error)
	GetCardById(userID, listID, itemID int) ([]entities.Cards, error)
	DeleteCardById(UserID, itemID int) error
}

type Service struct {
	Authenthication
	WordsList
	ListCards
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authenthication: NewAuthService(repos.Authenthication),
		WordsList:       NewWordsLists(repos.WordsList),
		ListCards:       NewCardsService(repos.ListCards),
	}
}
