package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"cnep-backend/source/models"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

type Client struct {
	Conn   *websocket.Conn
	UserID uint
}

type Message struct {
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	db         *gorm.DB
	mu         sync.RWMutex
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				client.Conn.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

func (h *Hub) handleMessage(message *Message) {
	// Forward message to receiver
	h.mu.RLock()
	if client, ok := h.clients[message.ReceiverID]; ok {
		err := client.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
		}
	}
	h.mu.RUnlock()

	// Asynchronously save message to database
	go func() {
		dbMessage := models.Message{
			SenderID:   message.SenderID,
			ReceiverID: message.ReceiverID,
			Content:    message.Content,
		}
		if err := h.db.Create(&dbMessage).Error; err != nil {
			log.Printf("Error saving message to database: %v", err)
		}
	}()
}

func (h *Hub) HandleWebSocket(c *websocket.Conn) {
	userID := c.Locals("userID").(uint)
	client := &Client{Conn: c, UserID: userID}

	h.register <- client

	defer func() {
		h.unregister <- client
	}()

	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			var message Message
			if err := json.Unmarshal(p, &message); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			message.SenderID = userID
			h.broadcast <- &message
		}
	}
}
