package repository

import (
	"ToDo/pkg/entities"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CardPostgres struct {
	db *sqlx.DB
}

func NewCardPostgres(db *sqlx.DB) *CardPostgres {
	return &CardPostgres{db: db}
}

func (c *CardPostgres) CreateCard(listID int, cards entities.Cards) (int, error) {
	var id int
	tx, err := c.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO cards (front, back, imagelink, voicemessage) VALUES ($1, $2, $3, $4) RETURNING id")
	row := tx.QueryRow(query, cards.Front, cards.Back, cards.ImageLink, cards.VoiceMessage)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO listCards (item_id, list_id) VALUES ($1, $2)")
	_, err = tx.Exec(query, id, listID)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CardPostgres) GetCardsInList(userID, listID int) ([]entities.Cards, error) {
	var response []entities.Cards
	query := fmt.Sprintf("SELECT cards.id, cards.front, cards.back, cards.imagelink, cards.voicemessage \nFROM cards\nINNER JOIN listcards lc ON lc.item_id = cards.id\nINNER JOIN userlists ul ON lc.list_id = ul.list_id\nWHERE ul.user_id = $1 AND ul.list_id = $2")
	if err := c.db.Select(&response, query, userID, listID); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CardPostgres) GetCardById(userID, listID, itemID int) ([]entities.Cards, error) {
	var response []entities.Cards
	query := fmt.Sprintf("SELECT cards.id, cards.front, cards.back, cards.imagelink, cards.voicemessage FROM cards INNER JOIN listcards lc ON lc.item_id = cards.id INNER JOIN userlists ul ON lc.list_id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2 AND lc.item_id=$3")
	if err := c.db.Select(&response, query, userID, listID, itemID); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CardPostgres) DeleteCardById(userID, itemID int) error {
	query := fmt.Sprintf("DELETE FROM cards USING listcards, userlists WHERE cards.id = listcards.item_id AND listcards.list_id = userlists.list_id AND userlists.user_id = $1 AND cards.id = $2")
	_, err := c.db.Exec(query, userID, itemID)
	if err != nil {
		return err
	}
	return nil
}
