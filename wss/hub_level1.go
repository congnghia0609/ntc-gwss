package wss

import (
	"bytes"
	"encoding/json"
	"log"
	"ntc-gwss/util"
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

func (h *HubLevel1) GetSizeClientLevel1(symbol string) int {
	return len(h.clients[symbol])
}

func (h *HubLevel1) BroadcastMsg(msg string) {
	util.TCF{
		Try: func() {
			if len(msg) > 0 {
				// log.Printf("message: %s", msg)

				message := []byte(msg)
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

				var data map[string]interface{}
				json.Unmarshal(message, &data)
				if data["s"] != nil {
					symbol := data["s"].(string)
					// log.Printf("HubLevel1.BroadcastMsg.symbol=%s", symbol)
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
		},
		Catch: func(e util.Exception) {
			log.Printf("HubLevel1.BroadcastMsg Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubLevel1) BroadcastMsgByte(message []byte) {
	util.TCF{
		Try: func() {
			if len(message) > 0 {
				// log.Printf("message: %s", message)
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

				var data map[string]interface{}
				json.Unmarshal(message, &data)
				if data["s"] != nil {
					symbol := data["s"].(string)
					// log.Printf("HubLevel1.BroadcastMsgByte.symbol=%s", symbol)
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
		},
		Catch: func(e util.Exception) {
			log.Printf("HubLevel1.BroadcastMsgByte Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubLevel1) run() {
	for {
		select {
		case client := <-h.register:
			util.TCF{
				Try: func() {
					if client != nil {
						if _, ok := h.clients[client.symbol]; ok {
							h.clients[client.symbol][client] = true
						} else {
							h.clients[client.symbol] = make(map[*ClientLevel1]bool)
							h.clients[client.symbol][client] = true
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel1.register Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		case client := <-h.unregister:
			util.TCF{
				Try: func() {
					if client != nil {
						if _, ok := h.clients[client.symbol][client]; ok {
							delete(h.clients[client.symbol], client)
							close(client.send)
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel1.unregister Caught %v\n", e)
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
						var data map[string]interface{}
						json.Unmarshal([]byte(message), &data)
						if data["s"] != nil {
							symbol := data["s"].(string)
							// log.Printf("HubLevel1.broadcast.symbol=%s", symbol)
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
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel1.broadcast Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		}
	}
}
