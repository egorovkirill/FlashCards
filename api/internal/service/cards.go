package service

import (
	"api/internal/repository/postgres"
	"api/pkg/entities"
)

type CardsService struct {
	repo postgres.ListCards
}

func NewCardsService(repo postgres.ListCards) *CardsService {

	return &CardsService{repo: repo}
}

func (r *CardsService) CreateCard(listID int, cards entities.Cards) error {
	return r.repo.CreateCard(listID, cards)
}

func (r *CardsService) GetCardsInList(userId, listID int) ([]entities.Cards, error) {
	return r.repo.GetCardsInList(userId, listID)
}

func (r *CardsService) GetCardById(userID, listID, itemID int) ([]entities.Cards, error) {
	return r.repo.GetCardById(userID, listID, itemID)
}

func (r *CardsService) DeleteCardById(UserID, itemID int) error {
	return r.repo.DeleteCardById(UserID, itemID)
}

func (r *CardsService) SetImageToCard(cardID int, image string) error {
	return r.repo.SetImageToCard(cardID, image)
}

func (r *CardsService) SetBackToCard(cardID int, translate string) error {
	return r.repo.SetImageToCard(cardID, translate)
}
