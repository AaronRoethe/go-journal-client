package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/AaronRoethe/go-journal-client/pocket"
)

type JournalEntry struct {
	JournalType string `json:"journalType"`
	Emotion     string `json:"emotion"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Author      string `json:"author"`
}

func AssembleMessage(buf *bytes.Buffer, tmplText string, answers JournalEntry) error {
	tmpl, err := template.New("").Parse(tmplText)
	if err != nil {
		return err
	}

	return tmpl.Execute(buf, answers)
}

func OutputStringTemplate(answers JournalEntry, tmplText string) ([]byte, error) {
	var buf bytes.Buffer
	if err := AssembleMessage(&buf, tmplText, answers); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Post_msg(answers JournalEntry) error {
	user, err := pocket.LoadUser()
	if err != nil {
		return fmt.Errorf("failed to load user details: %v", err)
	}

	reqBody := JournalEntry{
		JournalType: answers.JournalType,
		Emotion:     answers.Emotion,
		Title:       answers.Title,
		Body:        answers.Body,
		Author:      user.ID,
	}

	payloadBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal JournalEntry payload: %v", err)
	}
	domain := "go-journal.roethe.dev"
	if value, ok := os.LookupEnv("DOMAIN"); ok {
		domain = value
	}
	url := fmt.Sprintf("https://%s/api/collections/messages/records", domain)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %v", err)
	}

	// Add the authentication token to the request headers
	req.Header.Set("Authorization", user.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh authentication token: %v", resp.Status)
	}

	return nil
}

func Journal() {
	form, tmplText, err := LoadForm()
	if err != nil {
		log.Fatalf("failed to fill load form: %v", err)
	}

	answers, err := journalPrompt(form)
	if err != nil {
		log.Fatalf("failed to fill out form: %v", err)
	}

	err = Post_msg(answers)
	if err != nil {
		fmt.Printf("Error posting JournalEntry: %v\n", err)
	}

	output, err := OutputStringTemplate(answers, tmplText)
	if err != nil {
		log.Fatalf("failed to output string template: %v", err)
	}

	fmt.Println("Output:", string(output))
}
