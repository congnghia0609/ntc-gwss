/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"bytes"
	"encoding/json"
	"log"
	"ntc-gwss/util"
)

// HubLevel2 maintain the set of active clients and broadcasts message to the client.
type HubLevel2 struct {
	// Registered clients.
	clients map[string]map[*ClientLevel2]bool

	// Inbound message from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *ClientLevel2

	// Unregister requests from the clients.
	unregister chan *ClientLevel2
}

func newHubLevel2() *HubLevel2 {
	return &HubLevel2{
		broadcast:  make(chan []byte),
		register:   make(chan *ClientLevel2),
		unregister: make(chan *ClientLevel2),
		clients:    make(map[string]map[*ClientLevel2]bool),
	}
}

func (h *HubLevel2) GetSizeClientLevel2(key string) int {
	return len(h.clients[key])
}

func (h *HubLevel2) BroadcastMsg(msg string) {
	util.TCF{
		Try: func() {
			if len(msg) > 0 {
				// log.Printf("message: %s", msg)
				message := []byte(msg)
				h.BroadcastMsgByte(message)
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("HubLevel2.BroadcastMsg Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubLevel2) BroadcastMsgByte(message []byte) {
	util.TCF{
		Try: func() {
			if len(message) > 0 {
				// log.Printf("message: %s", message)
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				if len(message) > 0 {
					h.broadcast <- message
				}

				//var data map[string]interface{}
				//json.Unmarshal(message, &data)
				//if data["s"] != nil && data["tt"] != nil {
				//	symbol := data["s"].(string)
				//	tt := data["tt"].(string)
				//	// log.Printf("HubLevel2.BroadcastMsgByte {symbol=%s,typeTime=%s}", symbol, tt)
				//	if len(symbol) > 0 && len(tt) > 0 {
				//		key := symbol + "_" + tt
				//		for client := range h.clients[key] {
				//			select {
				//			case client.send <- message:
				//			default:
				//				close(client.send)
				//				delete(h.clients[key], client)
				//			}
				//		}
				//	}
				//}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("HubLevel2.BroadcastMsgByte Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubLevel2) run() {
	for {
		select {
		case client := <-h.register:
			util.TCF{
				Try: func() {
					if client != nil {
						key := client.symbol + "_" + client.typeTime
						if _, ok := h.clients[key]; ok {
							h.clients[key][client] = true
						} else {
							h.clients[key] = make(map[*ClientLevel2]bool)
							h.clients[key][client] = true
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel2.register Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		case client := <-h.unregister:
			util.TCF{
				Try: func() {
					if client != nil {
						key := client.symbol + "_" + client.typeTime
						if _, ok := h.clients[key][client]; ok {
							delete(h.clients[key], client)
							close(client.send)
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel2.unregister Caught %v\n", e)
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
						if data["s"] != nil && data["tt"] != nil {
							symbol := data["s"].(string)
							tt := data["tt"].(string)
							// log.Printf("HubLevel2.broadcast {symbol=%s,typeTime=%s}", symbol, tt)
							if len(symbol) > 0 && len(tt) > 0 {
								key := symbol + "_" + tt
								for client := range h.clients[key] {
									select {
									case client.send <- message:
									default:
										close(client.send)
										delete(h.clients[key], client)
									}
								}
							}
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubLevel2.broadcast Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		}
	}
}
