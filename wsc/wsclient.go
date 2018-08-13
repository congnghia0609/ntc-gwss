package wsc

import (
	"log"
	"net/url"
	"ntc-gwsc/util"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type UWSClient struct {
	name      string
	interrupt chan os.Signal
	done      chan struct{}
	url       url.URL
	conn      *websocket.Conn
}

var mapInstanceWSC = make(map[string]*UWSClient)

// GetInstanceWSC Singleton UWSClient
func GetInstanceWSC(name string) *UWSClient {
	log.Println("=================== UWSClient.GetInstanceWSC ===================")
	return mapInstanceWSC[name]
}

// NewInstanceWSC new object UWSClient
func NewInstanceWSC(name string, scheme string, host string, path string) (*UWSClient, error) {
	var uws *UWSClient
	var err error
	util.TCF{
		Try: func() {
			log.Println("=================== UWSClient.NewInstanceWSC ===================")
			log.Printf("+++++++++ name: %s", name)
			if len(name) <= 0 {
				name = uuid.Must(uuid.NewV4()).String()
			}
			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)
			done := make(chan struct{})
			url := url.URL{Scheme: scheme, Host: host, Path: path}
			log.Printf("connecting to %s", url.String())
			conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

			uws = &UWSClient{name: name, interrupt: interrupt, done: done, url: url, conn: conn}

			if err != nil {
				log.Println("[xxxxxx] ERROR:", err)
				//log.Fatal("dial:", err)
			}
			mapInstanceWSC[name] = uws
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.NewInstanceWSC Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
	return uws, err
}

func (uws *UWSClient) Close() {
	if uws.conn != nil {
		uws.conn.Close()
	}
}

// Reconnect is method auto-reconnect to websocket server.
func (uws *UWSClient) Reconnect() {
	// create ticket time for every 3 seconds
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	var count = 1
	for _ = range ticker.C {
		conn, _, err := websocket.DefaultDialer.Dial(uws.url.String(), nil)

		// Create new websocket
		log.Printf("\nRetry Connect [%s]: %d times\n", uws.url.String(), count)
		count = count + 1
		if err != nil {
			log.Printf("Dial failed [%s]: %s\n\n", uws.url.String(), err.Error())
		} else {
			uws.conn = conn
			break
		}
	}
}

// // Reconnect is method auto-reconnect to websocket server.
// func Reconnect(name string, scheme string, host string, path string) *UWSClient {
// 	// create ticket time for every 3 seconds
// 	ticker := time.NewTicker(time.Duration(3) * time.Second)
// 	var count = 1
// 	for _ = range ticker.C {
// 		var uws *UWSClient
// 		var err error
// 		uws, err = NewInstanceWSC(name, scheme, host, path)

// 		// Create new websocket
// 		fmt.Printf("\nRetry Connect : %d times\n", count)
// 		count = count + 1
// 		if err != nil {
// 			log.Printf("Dial failed: %s\n\n", err.Error())
// 		} else {
// 			return uws
// 		}
// 	}
// 	return nil
// }

// func main() {
// 	var uws *UWSClient
// 	// var err error
// 	uws = GetInstanceUWS()
// 	defer uws.Close()

// 	// Thread receive message.
// 	go func() {
// 		defer uws.Close()
// 		defer close(uws.done)
// 		for {
// 			_, message, err := uws.conn.ReadMessage()
// 			if err != nil {
// 				log.Println("read:", err)
// 				return
// 			}
// 			log.Printf("recv: %s", message)
// 		}
// 	}()

// 	// Thread send message.
// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case t := <-ticker.C:
// 			err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
// 			if err != nil {
// 				log.Println("write:", err)
// 				return
// 			}
// 		case <-uws.interrupt:
// 			log.Println("interrupt")
// 			// To cleanly close a connection, a client should send a close
// 			// frame and wait for the server to close the connection.
// 			err := uws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 			if err != nil {
// 				log.Println("write close:", err)
// 				return
// 			}
// 			select {
// 			case <-uws.done:
// 			case <-time.After(time.Second):
// 			}
// 			uws.Close()
// 			return
// 		}
// 	}
// }
