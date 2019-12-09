/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wsc

import (
	"fmt"
	"log"
	"ntc-gwss/util"
	"ntc-gwss/wss"
	"time"
	"ntc-gwss/conf"

	"github.com/gorilla/websocket"
)

func (wsc *UWSClient) recvCS() {
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

func (wsc *UWSClient) sendCS() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
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

func NewCSWSClient() *UWSClient {
	var cswsc *UWSClient
	c := conf.GetConfig()
	address := c.GetString("dataws.host") + ":" + c.GetString("dataws.port")
	log.Printf("################ CSWSClient[%s] start...", NameCSWSC)
	// ws://e-internal-data1:15401/dataws/stock
	cswsc, _ = NewInstanceWSC(NameCSWSC, "ws", address, "/dataws/stock")
	// cswsc, _ = NewInstanceWSC(NameCSWSC, "ws", "localhost:15601", "/ws/v1/cs/ETH_BTC@1h")
	//wss://engine2.kryptono.exchange/ws/v1/cs/ETH_BTC@1m
	// cswsc, _ = NewInstanceWSC(NameCSWSC, "wss", "engine2.kryptono.exchange", "/ws/v1/cs/ETH_BTC@1m")
	return cswsc
}

func (cswsc *UWSClient) StartCSWSClient() {
	// Thread receive message.
	go cswsc.recvCS()
	// Thread send message.
	//go cswsc.sendCS()
}
