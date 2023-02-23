package DALL_E_API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

type Request struct {
	Prompt         string `json:"prompt"`
	NumImages      int    `json:"num_images"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
}

var reqHeaders = map[string]string{
	"Content-Type":  "application/json",
	"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("API_SECRET")),
}

var apiURL = "https://api.openai.com/v1/images/generations"

func GenerateImage(prompt string) (string, error) {
	request := Request{
		Prompt:         prompt,
		NumImages:      1,
		Size:           "256x256",
		ResponseFormat: "url",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	for key, value := range reqHeaders {
		req.Header.Set(key, value)
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Data[0].URL, nil
}
