package KafkaRPC

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const API = "https://api.openai.com/v1/images/generations"

type LookupResponse struct {
	Head Head  `json:"head"`
	Def  []Def `json:"def"`
}

type Head struct{}

type Def struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Ts   string `json:"ts"`
	Tr   []Tr   `json:"tr"`
	Mean []Mean `json:"mean"`
	Syn  []Syn  `json:"syn"`
}

type Tr struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Syn  []Syn  `json:"syn"`
	Mean []Mean `json:"mean"`
	Asp  string `json:"asp"`
}

type Syn struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Gen  string `json:"gen"`
	Fr   int    `json:"fr"`
}

type Mean struct {
	Text string `json:"text"`
}

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

func GetImage(prompt string) (*http.Request, error) {
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
