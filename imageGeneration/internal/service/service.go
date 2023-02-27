package service

import (
	"imageGeneration/internal/database"
)

type ListCards interface {
	SetImageToCard(cardID int, prompt string) error
	SetTranslateToCard(cardID int, prompt string) error
}

type Service struct {
	ListCards
}

func NewService(repo *database.Repository) *Service {
	return &Service{
		ListCards: NewCreateCard(repo.ListCards),
	}
}
