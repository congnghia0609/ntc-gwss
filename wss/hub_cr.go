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
	"strings"
	"ntc-gwss/util"
)

// HubCR maintain the set of active clients and broadcasts message to the client.
type HubCR struct {
	// Registered clients Ticker24h.
	mapclient map[*ClientCR]bool
	// Registered clients DepthPrice & HistoryTrade.
	mapsymbolclient map[string]map[*ClientCR]bool

	// Inbound message from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *ClientCR

	// Unregister requests from the clients.
	unregister chan *ClientCR
}

func newHubCR() *HubCR {
	return &HubCR{
		broadcast:       make(chan []byte),
		register:        make(chan *ClientCR),
		unregister:      make(chan *ClientCR),
		mapclient:       make(map[*ClientCR]bool),
		mapsymbolclient: make(map[string]map[*ClientCR]bool),
	}
}

func (h *HubCR) GetSizeClientCR() int {
	return len(h.mapclient)
}

func (h *HubCR) GetSizeSymbolClientCR(symbol string) int {
	return len(h.mapsymbolclient[symbol])
}

func (h *HubCR) BroadcastMsg(msg string) {
	util.TCF{
		Try: func() {
			if len(msg) > 0 {
				// log.Printf("message: %s", msg)
				message := []byte(msg)
				h.BroadcastMsgByte(message)
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("HubCR.BroadcastMsg Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubCR) BroadcastMsgByte(message []byte) {
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
				//if data["et"] != nil {
				//	et := data["et"].(string)
				//	// broadcast Ticker24h
				//	if strings.EqualFold("tk", et) {
				//		TKDataCacheCR = string(message[:])
				//		for client := range h.mapclient {
				//			select {
				//			case client.send <- message:
				//			default:
				//				close(client.send)
				//				delete(h.mapclient, client)
				//			}
				//		}
				//	} else if strings.EqualFold("dp", et) || strings.EqualFold("ht", et) {
				//		// broadcast DepthPrice & HistoryTrade
				//		if data["s"] != nil {
				//			symbol := data["s"].(string)
				//			// log.Printf("HubCR.broadcast.symbol=%s", symbol)
				//			if len(symbol) > 0 {
				//				for client := range h.mapsymbolclient[symbol] {
				//					select {
				//					case client.send <- message:
				//					default:
				//						close(client.send)
				//						delete(h.mapsymbolclient[symbol], client)
				//					}
				//				}
				//			}
				//		}
				//	}
				//}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("HubCR.BroadcastMsgByte Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (h *HubCR) run() {
	for {
		select {
		case client := <-h.register:
			util.TCF{
				Try: func() {
					if client != nil {
						// register Ticker24h
						h.mapclient[client] = true
						// register DepthPrice & HistoryTrade
						if _, ok := h.mapsymbolclient[client.symbol]; ok {
							h.mapsymbolclient[client.symbol][client] = true
						} else {
							h.mapsymbolclient[client.symbol] = make(map[*ClientCR]bool)
							h.mapsymbolclient[client.symbol][client] = true
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubCR.register Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		case client := <-h.unregister:
			util.TCF{
				Try: func() {
					if client != nil {
						// unregister Ticker24h
						if _, ok := h.mapclient[client]; ok {
							delete(h.mapclient, client)
						}
						// unregister DepthPrice & HistoryTrade
						if _, ok := h.mapsymbolclient[client.symbol][client]; ok {
							delete(h.mapsymbolclient[client.symbol], client)
						}
						close(client.send)
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubCR.unregister Caught %v\n", e)
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
						if data["et"] != nil {
							et := data["et"].(string)
							// broadcast Ticker24h
							if strings.EqualFold("tk", et) {
								TKDataCacheCR = string(message[:])
								for client := range h.mapclient {
									select {
									case client.send <- message:
									default:
										close(client.send)
										delete(h.mapclient, client)
									}
								}
							} else if strings.EqualFold("dp", et) || strings.EqualFold("ht", et) {
								// broadcast DepthPrice & HistoryTrade
								if data["s"] != nil {
									symbol := data["s"].(string)
									// log.Printf("HubCR.broadcast.symbol=%s", symbol)
									if len(symbol) > 0 {
										for client := range h.mapsymbolclient[symbol] {
											select {
											case client.send <- message:
											default:
												close(client.send)
												delete(h.mapsymbolclient[symbol], client)
											}
										}
									}
								}
							}
						}
					}
				},
				Catch: func(e util.Exception) {
					log.Printf("HubCR.broadcast Caught %v\n", e)
				},
				Finally: func() {
					//log.Println("Finally...")
				},
			}.Do()
		}
	}
}
