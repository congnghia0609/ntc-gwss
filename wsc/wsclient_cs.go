/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wsc

import (
	"fmt"
	"log"
	"ntc-gwss/conf"
	"ntc-gwss/util"
	"ntc-gwss/wss"
	"time"

	"github.com/gorilla/websocket"
)

func (wsc *NWSClient) recvCS() {
	util.TCF{
		Try: func() {
			defer wsc.Close()
			defer close(wsc.done)
			for {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					wsc.Reconnect()
					// return
				}
				// log.Printf("recvCS: %s", message)
				if len(message) > 0 {
					// CSWSServer
					csws := wss.GetInstanceCS(wss.NameCSWSS)
					if csws != nil {
						csws.GetHub().BroadcastMsgByte(message)
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvCS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *NWSClient) sendCS() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					//err := nws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
					msec := t.UnixNano() / 1000000
					///// 1. Candlesticks Data.
					data := `{"tt":"1h","s":"ETH_BTC","t":` + fmt.Sprint(msec) + `,"e":"kline","k":{"c":"0.00028022","t":1533715200000,"v":"905062.00000000","h":"0.00028252","l":"0.00027787","o":"0.00027919"}}`
					err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(data))
					if err != nil {
						log.Println("write:", err)
						//return
					}
				case <-wsc.interrupt:
					log.Println("interrupt")
					// To cleanly close a connection, a client should send a close
					// frame and wait for the server to close the connection.
					err := wsc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-wsc.done:
					case <-time.After(time.Second):
					}
					wsc.Close()
					return
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.sendCS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

// NewCSWSClient new instance CSWSClient of NWSClient
func NewCSWSClient() *NWSClient {
	var cswsc *NWSClient
	c := conf.GetConfig()
	scheme := c.GetString(NameCSWSC + ".wsc.scheme")
	address := c.GetString(NameCSWSC + ".wsc.host")
	path := c.GetString(NameCSWSC + ".wsc.path")
	log.Printf("################ CSWSClient[%s] start...", NameCSWSC)
	cswsc, _ = NewInstanceWSC(NameCSWSC, scheme, address, path)
	// cswsc, _ = NewInstanceWSC(NameCSWSC, "ws", address, "/dataws/stock")
	// cswsc, _ = NewInstanceWSC(NameCSWSC, "ws", "localhost:15601", "/ws/v1/cs/ETH_BTC@1h")
	// cswsc, _ = NewInstanceWSC(NameCSWSC, "wss", "engine2.kryptono.exchange", "/ws/v1/cs/ETH_BTC@1m")
	return cswsc
}

// StartCSWSClient start
func (wsc *NWSClient) StartCSWSClient() {
	// Thread receive message.
	go wsc.recvCS()
	// Thread send message.
	//go wsc.sendCS()
}
