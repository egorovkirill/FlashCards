package database

import "github.com/jmoiron/sqlx"

type ListCards interface {
	//SetTranslateToCard(cardID int, translate string) error
	SetImageToCard(cardID int, image string) error
}

type Repository struct {
	ListCards
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		ListCards: NewCardPostgres(db),
	}
}
