/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"github.com/congnghia0609/ntc-gwss/conf"
	"log"
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

	// Maximum message size allowed from peer. Default: 10MB
	maxMessageSize = 10485760

	// Buffered channel send of client. Default: 100KB
	sendBuffer = 102400
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  102400, // 100 KB
	WriteBufferSize: 102400, // 100 KB
}

// MapSymbol Const
var MapSymbol = make(map[string]string)

// TypeTime type time
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

const (
	// NameDPWSS depthprice
	NameDPWSS = "depthprice"
	// NameCSWSS candlesticks
	NameCSWSS = "candlesticks"
	// NameHTWSS historytrade
	NameHTWSS = "historytrade"
	// NameTKWSS ticker24h
	NameTKWSS = "ticker24h"
	// NameCRWSS cerberus
	NameCRWSS = "cerberus"
)

// InitMapSymbol init map symbol
func InitMapSymbol() {
	c := conf.GetConfig()
	listpair := c.GetString("market.listpair")
	log.Printf("=========== listpair: %s", listpair)

	//var listpair = "ETH_BTC;KNOW_BTC;KNOW_ETH"
	var arrSymbol = strings.Split(listpair, ";")
	// log.Printf("arrSymbol: ", arrSymbol)
	for i := range arrSymbol {
		symbol := arrSymbol[i]
		MapSymbol[symbol] = symbol
	}
	log.Printf("=========== MapSymbol: %v", MapSymbol)
}

// ReloadMapSymbol update map symbol
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
		log.Printf("=========== reloadMapSymbol.MapSymbol: %v", MapSymbol)
	}
}
