package wss

import (
	"bytes"
	"log"
	"ntc-gwss/util"
)

// Hub maintain the set of active clients and broadcasts message to the client.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound message from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from the clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) GetSizeClient() int {
	return len(h.clients)
}

func (h *Hub) BroadcastMsg(msg string) {
	util.TCF{
		Try: func() {
			if len(msg) > 0 {
				// log.Printf("message: %s", msg)

				message := []byte(msg)
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				if len(message) > 0 {
					for client := range h.clients {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("Hub.BroadcastMsg Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *Hub) BroadcastMsgByte(message []byte) {
	util.TCF{
		Try: func() {
			if len(message) > 0 {
				// log.Printf("message: %s", message)
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				if len(message) > 0 {
					for client := range h.clients {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("Hub.BroadcastMsgByte Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			util.TCF{
				Try: func() {
					if client != nil {
						h.clients[client] = true
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("Hub.register Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		case client := <-h.unregister:
			util.TCF{
				Try: func() {
					if client != nil {
						if _, ok := h.clients[client]; ok {
							delete(h.clients, client)
							close(client.send)
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("Hub.unregister Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		case message := <-h.broadcast:
			util.TCF{
				Try: func() {
					if len(message) > 0 {
						// log.Printf("message: %s", message)
						for client := range h.clients {
							select {
							case client.send <- message:
							default:
								close(client.send)
								delete(h.clients, client)
							}
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("Hub.broadcast Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		}
	}
}
