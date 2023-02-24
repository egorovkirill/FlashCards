package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CardPostgres struct {
	db *sqlx.DB
}

func NewCardPostgres(db *sqlx.DB) *CardPostgres {
	return &CardPostgres{db: db}
}

//func (c *CardPostgres) SetTranslateToCard(cardID int, translate string) error {
//	query := fmt.Sprintf("UPDATE cards SET back = $1 WHERE id = $2")
//	_, err := c.db.Exec(query, translate, cardID)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (c *CardPostgres) SetImageToCard(cardID int, image string) error {
	query := fmt.Sprintf("UPDATE cards SET imagelink = $1 WHERE id = $2")
	_, err := c.db.Exec(query, image, cardID)
	if err != nil {
		return err
	}
	return nil
}
