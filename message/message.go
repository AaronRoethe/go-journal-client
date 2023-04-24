package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type MessageData struct {
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}

func SendMessage() error {
	message, err := FillOutForm()
	if err != nil {
		return err
	}

	data := MessageData{
		Msg:       string(message),
		Timestamp: time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(data)

	url := os.Getenv("POST_URL")
	if url == "" {
		url = "http://localhost:7071/api/endpoint"
	}

	req, _ := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %v", resp.StatusCode)
	}

	return nil
}
