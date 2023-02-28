package postgres

import (
	"api/pkg/entities"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"
)

type CardPostgres struct {
	db         *sqlx.DB
	kafkaAsync sarama.AsyncProducer
}

func NewCardPostgres(db *sqlx.DB, kafkaAsync sarama.AsyncProducer) *CardPostgres {
	return &CardPostgres{db: db, kafkaAsync: kafkaAsync}
}

func (c *CardPostgres) CreateCard(listID int, cards entities.Cards) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO cards (front, back, imagelink) VALUES ($1, $2, $3) RETURNING id")
	row := tx.QueryRow(query, cards.Front, cards.Back, cards.ImageLink)
	if err := row.Scan(&cards.Id); err != nil {
		_ = tx.Rollback()
		return err
	}
	query = fmt.Sprintf("INSERT INTO listCards (item_id, list_id) VALUES ($1, $2)")
	_, err = tx.Exec(query, cards.Id, listID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	data, err := json.Marshal(cards)
	// Write card data to Kafka topic
	message := &sarama.ProducerMessage{
		Topic: "cards_topic",
		Value: sarama.StringEncoder(data),
	}
	c.kafkaAsync.Input() <- message
	select {
	case err := <-c.kafkaAsync.Errors():
		return err
	case _ = <-c.kafkaAsync.Successes():
		return nil
	}
}

func (c *CardPostgres) GetCardsInList(userID, listID int) ([]entities.Cards, error) {
	var response []entities.Cards
	query := fmt.Sprintf("SELECT cards.id, cards.front, cards.back, cards.imagelink \nFROM cards\nINNER JOIN listcards lc ON lc.item_id = cards.id\nINNER JOIN userlists ul ON lc.list_id = ul.list_id\nWHERE ul.user_id = $1 AND ul.list_id = $2")
	if err := c.db.Select(&response, query, userID, listID); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CardPostgres) GetCardById(userID, listID, itemID int) ([]entities.Cards, error) {
	var response []entities.Cards
	query := fmt.Sprintf("SELECT cards.id, cards.front, cards.back, cards.imagelink FROM cards INNER JOIN listcards lc ON lc.item_id = cards.id INNER JOIN userlists ul ON lc.list_id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2 AND lc.item_id=$3")
	if err := c.db.Select(&response, query, userID, listID, itemID); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CardPostgres) SetImageToCard(cardID int, image string) error {
	query := fmt.Sprintf("UPDATE cards SET imagelink = $1 WHERE id = $2")
	_, err := c.db.Exec(query, image, cardID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardPostgres) SetTranslateToCard(cardID int, translate string) error {
	query := fmt.Sprintf("UPDATE cards SET back = $1 WHERE id = $2")
	_, err := c.db.Exec(query, translate, cardID)
	if err != nil {
		return err
	}
	return nil
}

func (c *CardPostgres) DeleteCardById(userID, itemID int) error {
	query := fmt.Sprintf("DELETE FROM cards USING listcards, userlists WHERE cards.id = listcards.item_id AND listcards.list_id = userlists.list_id AND userlists.user_id = $1 AND cards.id = $2")
	_, err := c.db.Exec(query, userID, itemID)
	if err != nil {
		return err
	}
	return nil
}
