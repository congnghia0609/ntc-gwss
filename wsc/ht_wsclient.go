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

func (wsc *UWSClient) recvHT() {
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
				// log.Printf("recvHT: %s", message)
				if len(message) > 0 {
					htws := wss.GetInstanceHT(wss.NameHTWSS)
					if htws != nil {
						htws.GetHub().BroadcastMsgByte(message)
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvHT Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *UWSClient) sendHT() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))

					msec := t.UnixNano() / 1000000

					///// 1. Historytrade Data.
					data := `{"p":"0.05567000","q":"1.84100000","c":1533886283334,"s":"ETH_BTC","t":` + fmt.Sprint(msec) + `,"e":"history_trade","k":514102,"m":true}`

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
			log.Printf("wsc.sendHT Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func NewHTWSClient() *UWSClient {
	var htwsc *UWSClient
	c := conf.GetConfig()
	address := c.GetString("dataws.host") + ":" + c.GetString("dataws.port")
	log.Printf("################ HTWSClient[%s] start...", NameHTWSC)
	// ws://e-internal-data1:15401/dataws/history
	htwsc, _ = NewInstanceWSC(NameHTWSC, "ws", address, "/dataws/history")
	// htwsc, _ = NewInstanceWSC(NameHTWSC, "ws", "localhost:15701", "/ws/v1/ht/ETH_BTC")
	// wss://engine2.kryptono.exchange/ws/v1/ht/ETH_BTC
	// htwsc, _ = NewInstanceWSC(NameHTWSC, "wss", "engine2.kryptono.exchange", "/ws/v1/ht/ETH_BTC")
	return htwsc
}

func (htwsc *UWSClient) StartHTWSClient() {
	// Thread receive message.
	go htwsc.recvHT()
	// Thread send message.
	//go htwsc.sendHT()
}
