package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gorilla/websocket"
	"github.com/user/vigilante/internal/storage"
)

// Hub maintains active WebSocket connections.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan storage.Anomaly
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

// NewHub initiates an alerting hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan storage.Anomaly),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

// Run holds the multiplexing event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case anomaly := <-h.broadcast:
			for client := range h.clients {
				err := client.WriteJSON(anomaly)
				if err != nil {
					client.Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

// SendSlackWebhook posts an anomaly to a Slack endpoint.
func SendSlackWebhook(ctx context.Context, webhookURL string, anomaly storage.Anomaly) error {
	payload := map[string]string{
		"text": "🚨 *New Anomaly Detected*\nType: " + anomaly.AnomalyType + "\n" + anomaly.Description,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// SendEmail dispatches an anomaly via SMTP.
func SendEmail(smtpHost, smtpPort, user, pass, to string, anomaly storage.Anomaly) error {
	auth := smtp.PlainAuth("", user, pass, smtpHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Vigilante Alert: " + anomaly.AnomalyType + "\r\n" +
		"\r\n" +
		anomaly.Description + "\r\n" +
		"- " + anomaly.SuggestedFix)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, user, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}
	return err
}
