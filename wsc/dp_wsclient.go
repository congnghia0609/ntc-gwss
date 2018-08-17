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

func (wsc *UWSClient) recvDP() {
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
				// log.Printf("recvDP: %s", message)
				if len(message) > 0 {
					// DPWSServer
					dpws := wss.GetInstanceDP(wss.NameDPWSS)
					if dpws != nil {
						dpws.GetHub().BroadcastMsgByte(message)
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvDP Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *UWSClient) sendDP() {
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
					data := `{"a":[],"b":[["379.11400000", "0.03203000"]],"s":"ETH_BTC","t":"` + fmt.Sprint(msec) + `","e":"depthUpdate"}`

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
			log.Printf("wsc.sendDP Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func NewDPWSClient() *UWSClient {
	var dpwsc *UWSClient
	c := conf.GetConfig()
	address := c.GetString("dataws.host") + ":" + c.GetString("dataws.port")
	log.Printf("################ DPWSClient[%s] start...", NameDPWSC)
	// ws://e-internal-data1:15401/dataws/depth
	dpwsc, _ = NewInstanceWSC(NameDPWSC, "ws", address, "/dataws/depth")
	// dpwsc, _ = NewInstanceWSC(NameDPWSC, "ws", "localhost:15501", "/ws/v1/dp/ETH_BTC")
	//wss://engine2.kryptono.exchange/ws/v1/dp/ETH_BTC
	// dpwsc, _ = NewInstanceWSC(NameDPWSC, "wss", "engine2.kryptono.exchange", "/ws/v1/dp/ETH_BTC")
	return dpwsc
}

func (dpwsc *UWSClient) StartDPWSClient() {
	// Thread receive message.
	go dpwsc.recvDP()
	// Thread send message.
	//go dpwsc.sendDP()
}
