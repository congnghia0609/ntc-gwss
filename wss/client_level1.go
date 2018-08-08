package wss

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"ntc-gwss/util"
	"time"

	"github.com/gorilla/websocket"
)

// ClientLevel1 is a middleman between the websocket connection and the hub.
type ClientLevel1 struct {
	// symbol
	symbol string

	// HubLevel1
	hub *HubLevel1

	// The websocket connnection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application ensures that
// there is at most one reader on a connection by execution by executing all reads from this goroutine.
func (c *ClientLevel1) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// Switch Case: Process Business. Here simple broadcast message to all client in hub.
		c.hub.broadcast <- message
	}
}

// writePump pumps message from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The application ensures that
// there is at most one write to a connection by executing all writes from this goroutine.
func (c *ClientLevel1) writePump() {
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
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

// respMsg is send message to the current websocket message.
func (c *ClientLevel1) respMsg(message string) {
	util.TCF{
		Try: func() {
			if len(message) > 0 {
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write([]byte(message))

				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write(newline)
					w.Write(<-c.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			}
		},
		Catch: func(e util.Exception) {
			fmt.Printf("ClientLevel1.respMsg Caught %v\n", e)
		},
		Finally: func() {
			//fmt.Println("Finally...")
		},
	}.Do()
}

// serveWsLevel1 handles websocket requests from the peer.
func serveWsLevel1(symbol string, hub *HubLevel1, w http.ResponseWriter, r *http.Request) {
	util.TCF{
		Try: func() {
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			client := &ClientLevel1{symbol: symbol, hub: hub, conn: conn, send: make(chan []byte, 256)}
			client.hub.register <- client

			// Allow collection of memory referenced by the caller by doing all work in new goroutines.
			go client.writePump()
			go client.readPump()

			// Push message connected successfully.
			msgsc := `{"err":0,"msg":"Connected sucessfully"}`
			log.Println(msgsc)
			client.respMsg(msgsc)
		},
		Catch: func(e util.Exception) {
			fmt.Printf("serveWsLevel1 Caught %v\n", e)
		},
		Finally: func() {
			//fmt.Println("Finally...")
		},
	}.Do()
}