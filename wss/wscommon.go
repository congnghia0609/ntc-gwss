package wss

import (
	"log"
	"ntc-gwss/conf"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

//// Declare Global
// Websocket Const
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  102400,
	WriteBufferSize: 102400,
}

// MapSymbol Const
var MapSymbol = make(map[string]string)
var TypeTime = map[string]string{
	"1m":  "1m",
	"5m":  "5m",
	"15m": "15m",
	"30m": "30m",
	"1h":  "1h",
	"2h":  "2h",
	"4h":  "4h",
	"6h":  "6h",
	"12h": "12h",
	"1d":  "1d",
	"1w":  "1w",
}

func InitMapSymbol() {
	c := config.GetConfig()
	listpair := c.GetString("market.listpair")
	log.Printf("=========== listpair: %s", listpair)

	//var listpair = "ETH_BTC;KNOW_BTC;KNOW_ETH"
	var arrSymbol = strings.Split(listpair, ";")
	// log.Printf("arrSymbol: ", arrSymbol)
	for i := range arrSymbol {
		symbol := arrSymbol[i]
		MapSymbol[symbol] = symbol
	}
	log.Printf("=========== MapSymbol: ", MapSymbol)
}

func ReloadMapSymbol(listpair string) {
	log.Printf("=========== reloadMapSymbol.listpair: %s", listpair)
	if listpair != "" {
		var arrSymbol = strings.Split(listpair, ";")
		// log.Printf("arrSymbol: ", arrSymbol)
		for i := range arrSymbol {
			symbol := arrSymbol[i]
			//// If not exist, add to MapSymbol
			if _, ok := MapSymbol[symbol]; !ok {
				MapSymbol[symbol] = symbol
			}
		}
		log.Printf("=========== reloadMapSymbol.MapSymbol: ", MapSymbol)
	}
}
