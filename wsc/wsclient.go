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
