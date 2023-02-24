package service

import (
	"encoding/json"
	"imageGeneration/internal/database"
	"io"
	"net/http"
)

type CreateCard struct {
	repo database.ListCards
}

func NewCreateCard(repo database.ListCards) *CreateCard {

	return &CreateCard{repo: repo}
}

func (c *CreateCard) SetImageToCard(cardID int, prompt string) error {
	req, err := createRequest(prompt)
	if err != nil {
		return err
	}

	headers := createHeaders()
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	err = c.repo.SetImageToCard(cardID, response.Data[0].URL)
	if err != nil {
		return err
	}
	return nil

}
