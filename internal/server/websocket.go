package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var mu sync.Mutex

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// WebSocketHandler встановлює з'єднання з клієнтом
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading WebSocket:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = ""
	mu.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			log.Println("Error reading message:", err)
			break
		}

		mu.Lock()
		if clients[conn] == "" {
			clients[conn] = msg.Username
		} else {
			msg.Username = clients[conn]
			broadcast <- msg
		}
		mu.Unlock()
	}
}

// BroadcastMessages відправляє повідомлення всім клієнтам
func BroadcastMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

// InitWebSocket ініціалізує WebSocket
func InitWebSocket() {
	go BroadcastMessages()
}
