package database

import "github.com/jmoiron/sqlx"
import _ "github.com/lib/pq"

type ListCards interface {
	SetImageToCard(cardID int, image string) error
	SetTranslateToCard(cardID int, translate string) error
}

type Repository struct {
	ListCards
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		ListCards: NewCardPostgres(db),
	}
}
