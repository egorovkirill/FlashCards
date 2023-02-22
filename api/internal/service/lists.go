package service

import (
	"api/internal/repository"
	"api/pkg/entities"
)

type WordsLists struct {
	repo repository.WordsList
}

func NewWordsLists(repo repository.WordsList) *WordsLists {
	return &WordsLists{repo: repo}
}

func (r *WordsLists) CreateList(userId int, title string) (int, error) {
	return r.repo.CreateList(userId, title)
}

func (r *WordsLists) GetLists(userId int) ([]entities.Lists, error) {
	return r.repo.GetLists(userId)
}

func (r *WordsLists) GetListById(userId int, listId int) ([]entities.Lists, error) {
	return r.repo.GetListById(userId, listId)
}

func (r *WordsLists) UpdateListById(userId, ListId int, title string) error {
	return r.repo.UpdateListById(userId, ListId, title)
}

func (r *WordsLists) DeleteListById(userId, ListId int) error {
	return r.repo.DeleteListById(userId, ListId)
}
