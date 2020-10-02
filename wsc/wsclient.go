/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wsc

import (
	"log"
	"net/url"
	"ntc-gwss/util"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// NWSClient struct
type NWSClient struct {
	name      string
	interrupt chan os.Signal
	done      chan struct{}
	url       url.URL
	conn      *websocket.Conn
}

var mapInstanceWSC = make(map[string]*NWSClient)

// GetInstanceWSC Singleton NWSClient
func GetInstanceWSC(name string) *NWSClient {
	log.Println("=================== NWSClient.GetInstanceWSC ===================")
	return mapInstanceWSC[name]
}

// NewInstanceWSC new object NWSClient
func NewInstanceWSC(name string, scheme string, host string, path string) (*NWSClient, error) {
	var nws *NWSClient
	var err error
	util.TCF{
		Try: func() {
			log.Println("=================== NWSClient.NewInstanceWSC ===================")
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

			nws = &NWSClient{name: name, interrupt: interrupt, done: done, url: url, conn: conn}

			if err != nil {
				log.Println("[xxxxxx] ERROR:", err)
				//log.Fatal("dial:", err)
			}
			mapInstanceWSC[name] = nws
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.NewInstanceWSC Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
	return nws, err
}

// Close NWSClient
func (nws *NWSClient) Close() {
	if nws.conn != nil {
		nws.conn.Close()
	}
}

// Reconnect is method auto-reconnect to websocket server.
func (nws *NWSClient) Reconnect() {
	// create ticket time for every 3 seconds
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	var count = 1
	for range ticker.C {
		conn, _, err := websocket.DefaultDialer.Dial(nws.url.String(), nil)
		// Create new websocket
		log.Printf("\nRetry Connect [%s]: %d times\n", nws.url.String(), count)
		count = count + 1
		if err != nil {
			log.Printf("Dial failed [%s]: %s\n\n", nws.url.String(), err.Error())
		} else {
			nws.conn.Close()
			nws.conn = conn
			break
		}
	}
}
