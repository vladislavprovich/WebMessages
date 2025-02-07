package server

import (
	"messenger/internal/models"
	"net/http"
	"sync"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

const blockChain = 100

type WebSocket interface {
	WebSocketHandler(w http.ResponseWriter, r *http.Request)
	BroadcastMessages()
}

type webSocket struct {
	log       *zap.Logger
	clients   map[*websocket.Conn]string
	broadcast chan models.Message
	upgrader  websocket.Upgrader
	mu        sync.RWMutex
}

func NewWebSocket() WebSocket {
	ws := &webSocket{
		log:       zap.NewExample(),
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan models.Message, blockChain),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
	}
	// Run for message.
	go ws.BroadcastMessages()
	return ws
}

func (s *webSocket) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("websocket upgrade error", zap.Error(err))
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		s.log.Info("Client disconnected", zap.String("client", conn.RemoteAddr().String()))
		err = conn.Close()
		if err != nil {
			s.log.Error("close websocket connection error", zap.Error(err))
		}
	}()

	s.log.Info("Client connected", zap.String("client", conn.RemoteAddr().String()))

	s.mu.Lock()
	s.clients[conn] = ""
	s.mu.Unlock()

	for {
		var msg models.Message
		err = conn.ReadJSON(&msg)
		if err != nil {
			s.log.Error("websocket conn read error", zap.Error(err))
			break
		}

		s.mu.Lock()
		if s.clients[conn] == "" {
			s.clients[conn] = msg.Username
		} else {
			msg.Username = s.clients[conn]
			select {
			case s.broadcast <- msg:
			default:
				s.log.Warn("Broadcast channel is full, dropping message")
			}
		}
		s.mu.Unlock()
	}
}

func (s *webSocket) BroadcastMessages() {
	for msg := range s.broadcast {
		s.mu.RLock()
		for client := range s.clients {
			if err := client.WriteJSON(msg); err != nil {
				s.log.Error("Error sending message", zap.Error(err))
				client.Close()
				s.mu.RUnlock()
				s.mu.Lock()
				delete(s.clients, client)
				s.mu.Unlock()
				s.mu.RLock()
			}
		}
		s.mu.RUnlock()
	}
}
