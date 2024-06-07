package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Payload struct defines the structure of the incoming request body
type Payload struct {
	CreatorID string    `json:"creatorId"`
	CreatedTs time.Time `json:"createdTs"`
	Memo      string    `json:"memo"`
}

const telegramAPI = "https://api.telegram.org/bot%s/sendMessage"

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func sendMessage(botToken, chatID, text string) error {
	message := TelegramMessage{
		ChatID: chatID,
		Text:   text,
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %v", err)
	}

	url := fmt.Sprintf(telegramAPI, botToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}

// Handler function to process incoming requests
func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	botToken := "7471887116:AAGpw9WlsYDEwXa1elpf65cfVQeYwRf08YE"
	chatID := "-1002141513360"
	text := payload.CreatorID + "tarafından paylaşıldı:\n" + payload.Memo

	if err := sendMessage(botToken, chatID, text); err != nil {
		fmt.Printf("Could not send message: %v\n", err)
	} else {
		fmt.Println("Message sent successfully")
	}

	// Process the payload as needed (logging in this example)
	log.Printf("Received payload: %+v", payload)

	// Send a response back to the client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func main() {
	http.HandleFunc("/api/v1/webhook", handler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
