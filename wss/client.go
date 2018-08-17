/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"bytes"
	"log"
	"net/http"
	"ntc-gwss/util"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connnection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application ensures that
// there is at most one reader on a connection by execution by executing all reads from this goroutine.
func (c *Client) readPump() {
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
		if len(message) > 0 {
			message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		}
		// Switch Case: Process Business. Here simple broadcast message to all client in hub.
		// c.hub.broadcast <- message

		// Not receive message from Client. {"msg":"Message invalid","err":-1}
		msg := `{"err":-1,"msg":"Message invalid"}`
		c.respMsg(msg)
	}
}

// writePump pumps message from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The application ensures that
// there is at most one write to a connection by executing all writes from this goroutine.
func (c *Client) writePump() {
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
func (c *Client) respMsg(message string) {
	util.TCF{
		Try: func() {
			if len(message) > 0 {
				msg := []byte(message)
				c.send <- msg
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("Client.respMsg Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

// respMsgByte is send message to the current websocket message.
func (c *Client) respMsgByte(msg []byte) {
	util.TCF{
		Try: func() {
			if len(msg) > 0 {
				c.send <- msg
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("Client.respMsgByte Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	util.TCF{
		Try: func() {
			upgrader.CheckOrigin = func(r *http.Request) bool { return true }
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}

			client := &Client{hub: hub, conn: conn, send: make(chan []byte, sendBuffer)}
			client.hub.register <- client

			// Allow collection of memory referenced by the caller by doing all work in new goroutines.
			go client.writePump()
			go client.readPump()

			// Push message connected successfully.
			msgsc := `{"err":0,"msg":"Connected sucessfully"}`
			client.respMsg(msgsc)
			// log.Printf("==============>>>>>>>>>>>>>> TKDataCache: %s", TKDataCache)
			time.Sleep(100 * time.Millisecond)
			if len(TKDataCache) > 0 {
				client.respMsg(TKDataCache)
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("serveWs Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}
