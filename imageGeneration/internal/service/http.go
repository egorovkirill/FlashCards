package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const API = "https://api.openai.com/v1/images/generations"

type Request struct {
	Prompt         string `json:"prompt"`
	NumImages      int    `json:"num_images"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

type Response struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func createRequest(prompt string) (*http.Request, error) {
	request := Request{
		Prompt:         prompt,
		NumImages:      1,
		Size:           "512x512",
		ResponseFormat: "url",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", API, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	return req, nil
}

func createHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("API_SECRET")),
	}
}
