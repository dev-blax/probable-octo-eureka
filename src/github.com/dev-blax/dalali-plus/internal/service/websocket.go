package service

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketService struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mutex     sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

func (ws *WebSocketService) HandleConnections(c *websocket.Conn) {
	defer func() {
		ws.Mutex.Lock()
		delete(ws.Clients, c)
		ws.Mutex.Unlock()
		c.Close()
	}()

	ws.Mutex.Lock()
	ws.Clients[c] = true
	ws.Mutex.Unlock()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		ws.Broadcast <- msg
	}
}

func (ws *WebSocketService) HandleMessages() {
	for {
		msg := <-ws.Broadcast
		ws.Mutex.Lock()
		for client := range ws.Clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("error: %v ", err)
				client.Close()
				delete(ws.Clients, client)
			}
		}

		ws.Mutex.Unlock()
	}
}
