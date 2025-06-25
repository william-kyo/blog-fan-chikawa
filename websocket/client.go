package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"blog-fanchiikawa-service/service"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ID   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Type      string      `json:"type"`
	ChatID    int64       `json:"chatId,omitempty"`
	Content   string      `json:"content,omitempty"`
	MessageID string      `json:"messageId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		log.Printf("Received message: %+v", msg)
		
		// Handle different message types
		switch msg.Type {
		case "send_message":
			c.handleSendMessage(msg)
		case "ping":
			c.handlePing(msg)
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) SendMessage(message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
	default:
		close(c.send)
		delete(c.hub.clients, c)
	}

	return nil
}

func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func (c *Client) handleSendMessage(msg Message) {
	// Extract chatId from message
	chatId := msg.ChatID
	if chatId == 0 {
		c.sendErrorResponse(msg.MessageID, "Chat ID is required")
		return
	}

	// Extract message content
	content := msg.Content
	if content == "" {
		c.sendErrorResponse(msg.MessageID, "Message content is required")
		return
	}

	// Call chat service to send message
	ctx := context.Background()
	req := &service.SendMessageRequest{
		ChatID:  chatId,
		Message: content,
	}

	response, err := c.hub.chatService.SendMessage(ctx, req)
	if err != nil {
		c.sendErrorResponse(msg.MessageID, err.Error())
		return
	}

	// Send success response with bot message
	successMsg := &Message{
		Type:      "message_response",
		MessageID: msg.MessageID,
		Data:      response,
	}

	if err := c.SendMessage(successMsg); err != nil {
		log.Printf("Failed to send success response: %v", err)
	}
}

func (c *Client) handlePing(msg Message) {
	pongMsg := &Message{
		Type:      "pong",
		MessageID: msg.MessageID,
	}
	
	if err := c.SendMessage(pongMsg); err != nil {
		log.Printf("Failed to send pong response: %v", err)
	}
}

func (c *Client) sendErrorResponse(messageID, errorMsg string) {
	errMsg := &Message{
		Type:      "error",
		MessageID: messageID,
		Error:     errorMsg,
	}
	
	if err := c.SendMessage(errMsg); err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}