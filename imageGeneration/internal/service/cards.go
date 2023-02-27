package service

import (
	"imageGeneration/internal/database"
)

type CreateCard struct {
	repo database.ListCards
}

func NewCreateCard(repo database.ListCards) *CreateCard {

	return &CreateCard{repo: repo}
}

func (c *CreateCard) SetImageToCard(cardID int, prompt string) error {
	err := c.repo.SetImageToCard(cardID, prompt)
	if err != nil {
		return err
	}
	return nil
}

func (c *CreateCard) SetTranslateToCard(cardID int, prompt string) error {
	err := c.repo.SetTranslateToCard(cardID, prompt)
	if err != nil {
		return err
	}
	return nil
}
