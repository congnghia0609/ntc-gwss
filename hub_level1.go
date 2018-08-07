package main

import (
	"encoding/json"
	"log"
)

// HubLevel1 maintain the set of active clients and broadcasts message to the client.
type HubLevel1 struct {
	// Registered clients.
	clients map[string]map[*ClientLevel1]bool

	// Inbound message from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *ClientLevel1

	// Unregister requests from the clients.
	unregister chan *ClientLevel1
}

func newHubLevel1() *HubLevel1 {
	return &HubLevel1{
		broadcast:  make(chan []byte),
		register:   make(chan *ClientLevel1),
		unregister: make(chan *ClientLevel1),
		clients:    make(map[string]map[*ClientLevel1]bool),
	}
}

func (h *HubLevel1) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.clients[client.symbol]; ok {
				h.clients[client.symbol][client] = true
			} else {
				h.clients[client.symbol] = make(map[*ClientLevel1]bool)
				h.clients[client.symbol][client] = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.symbol][client]; ok {
				delete(h.clients[client.symbol], client)
				close(client.send)
			}
		case message := <-h.broadcast:
			if len(message) > 0 {
				log.Printf("message: %s", message)
				var data map[string]interface{}
				json.Unmarshal([]byte(message), &data)
				if data["s"] != nil {
					symbol := data["s"].(string)
					log.Printf("HubLevel1.broadcast.symbol=%s", symbol)
					if len(symbol) > 0 {
						for client := range h.clients[symbol] {
							select {
							case client.send <- message:
							default:
								close(client.send)
								delete(h.clients[symbol], client)
							}
						}
					}
				}
			}
		}
	}
}
