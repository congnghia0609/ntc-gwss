package wsc

import (
	"fmt"
	"log"
	"ntc-gwsc/util"
	"time"

	"github.com/gorilla/websocket"
)

func (wsc *UWSClient) recvRS() {
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
				//log.Printf("url: %s", wsc.url.String())
				log.Printf("recv: %s", message)
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvRS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *UWSClient) sendRS() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case t := <-ticker.C:
					//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))

					msec := t.UnixNano() / 1000000

					///// 1. Ticker24h Data.
					data := `{"t":` + fmt.Sprint(msec) + `,"list_symbol":"BTC_USDT;ETH_USDT;KNOW_USDT;GTO_USDT"}`
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
			log.Printf("wsc.sendRS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func NewRSWSClient() *UWSClient {
	var rswsc *UWSClient
	// var err error
	rswsc, _ = NewInstanceWSC(NameRSWSC, "ws", "localhost:15801", "/ws/v1/tk")
	//defer uws.Close()
	return rswsc
}

func (tkwsc *UWSClient) StartRSWSClient() {
	// Thread receive message.
	go tkwsc.recvRS()
	// Thread send message.
	go tkwsc.sendRS()
}
