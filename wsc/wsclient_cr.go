/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wsc

import (
	"fmt"
	"log"
	"time"
	"ntc-gwss/conf"
	"ntc-gwss/util"
	"ntc-gwss/wss"

	"github.com/gorilla/websocket"
)

func (wsc *UWSClient) recvCR() {
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
				// log.Printf("recvCR: %s", message)
				if len(message) > 0 {
					// CRWSServer
					crws := wss.GetInstanceCR(wss.NameCRWSS)
					if crws != nil {
						crws.GetHub().BroadcastMsgByte(message)
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvCR Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *UWSClient) sendCR() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case t := <-ticker.C:
					//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
					msec := t.UnixNano() / 1000000
					///// 1. DepthPrice Data.
					data := `{"et":"dp","s":"ETH_BTC",{"a":[],"b":[["379.11400000", "0.03203000"]],"s":"ETH_BTC","t":"` + fmt.Sprint(msec) + `","e":"depthUpdate"}}`
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
			log.Printf("wsc.sendCR Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func NewCRWSClient() *UWSClient {
	var crwsc *UWSClient
	c := conf.GetConfig()
	address := c.GetString("dataws.host") + ":" + c.GetString("dataws.port")
	log.Printf("################ CRWSClient[%s] start...", NameCRWSC)
	crwsc, _ = NewInstanceWSC(NameCRWSC, "ws", address, "/dataws/cerberus")
	// crwsc, _ = NewInstanceWSC(NameCRWSC, "ws", "localhost:15501", "/ws/v1/cr/ETH_BTC")
	// crwsc, _ = NewInstanceWSC(NameCRWSC, "wss", "engine2.kryptono.exchange", "/ws/v1/cr/ETH_BTC")
	return crwsc
}

func (crwsc *UWSClient) StartCRWSClient() {
	// Thread receive message.
	go crwsc.recvCR()
	// Thread send message.
	//go crwsc.sendCR()
}
