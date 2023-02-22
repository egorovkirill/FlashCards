package repository

import (
	"api/pkg/entities"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ListsPostgres struct {
	db *sqlx.DB
}

func NewCreateListPostgres(db *sqlx.DB) *ListsPostgres {
	return &ListsPostgres{db: db}
}

func (c *ListsPostgres) CreateList(userId int, title string) (int, error) {
	var id int
	tx, _ := c.db.Begin()
	query := fmt.Sprintf("INSERT INTO lists (title) VALUES ($1) RETURNING id")
	row := tx.QueryRow(query, title)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	query = fmt.Sprintf("INSERT INTO userlists (user_id, list_id) VALUES ($1, $2) RETURNING id")
	row = tx.QueryRow(query, userId, id)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	return id, nil
}

func (c *ListsPostgres) GetLists(userId int) ([]entities.Lists, error) {
	var response []entities.Lists
	query := fmt.Sprintf("SELECT lists.id, title FROM lists INNER JOIN userLists ON lists.id = userlists.list_id WHERE user_id = $1")
	if err := c.db.Select(&response, query, userId); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *ListsPostgres) GetListById(userId, ListId int) ([]entities.Lists, error) {
	var response []entities.Lists
	query := fmt.Sprintf("SELECT lists.id, title FROM lists INNER JOIN userLists ON lists.id = userlists.list_id WHERE user_id = $1 AND list_id = $2")
	if err := c.db.Select(&response, query, userId, ListId); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *ListsPostgres) UpdateListById(userId, ListId int, title string) error {
	query := fmt.Sprintf("UPDATE lists SET title = $1 WHERE id = (SELECT lists.id FROM lists INNER JOIN userLists ON lists.id = userlists.list_id  WHERE user_id = $2 AND list_id = $3)")
	_, err := c.db.Exec(query, title, userId, ListId)
	if err != nil {
		return err
	}
	return nil
}

func (c *ListsPostgres) DeleteListById(userId, ListId int) error {
	query := fmt.Sprintf("DELETE FROM lists WHERE id = (SELECT lists.id FROM lists INNER JOIN userLists ON lists.id = userlists.list_id  WHERE user_id = $1 AND list_id = $2)")
	_, err := c.db.Exec(query, userId, ListId)
	if err != nil {
		return err
	}
	return nil
}
