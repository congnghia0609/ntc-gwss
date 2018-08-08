package main

import (
	"flag"
	"fmt"
	"log"
	"ntc-gwss/conf"
	"ntc-gwss/wss"
	"os"
	"strings"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//// Declare Global
var dpwss *wss.DPWSServer
var cswss *wss.CSWSServer

func InitMapSymbol() {
	c := config.GetConfig()
	listpair := c.GetString("market.listpair")
	log.Printf("=========== listpair: %s", listpair)

	//var listpair = "ETH_BTC;KNOW_BTC;KNOW_ETH"
	var arrSymbol = strings.Split(listpair, ";")
	// log.Printf("arrSymbol: ", arrSymbol)
	for i := range arrSymbol {
		symbol := arrSymbol[i]
		wss.MapSymbol[symbol] = symbol
	}
	log.Printf("=========== MapSymbol: ", wss.MapSymbol)
}

func ReloadMapSymbol(listpair string) {
	log.Printf("=========== reloadMapSymbol.listpair: %s", listpair)
	if listpair != "" {
		var arrSymbol = strings.Split(listpair, ";")
		// log.Printf("arrSymbol: ", arrSymbol)
		for i := range arrSymbol {
			symbol := arrSymbol[i]
			//// If not exist, add to MapSymbol
			if _, ok := wss.MapSymbol[symbol]; !ok {
				wss.MapSymbol[symbol] = symbol
			}
		}
		log.Printf("=========== reloadMapSymbol.MapSymbol: ", wss.MapSymbol)
	}
}

func main() {
	finish := make(chan bool)

	//// init Configuration
	environment := flag.String("e", "development", "run project with mode [-e development | test | production]")
	flag.Usage = func() {
		fmt.Println("Usage: [appname] -e development|production")
		os.Exit(1)
	}
	flag.Parse()
	log.Printf("============== environment: %s", *environment)
	config.Init(*environment)

	//// initMapSymbol
	InitMapSymbol()

	// // NewWSServer
	// go wss.NewWSServer("ntc")

	//// Run DPWSServer
	dpwss = wss.NewDPWSServer("depthprice")
	log.Printf("======= DPWSServer[%s] is ready...", dpwss.GetName())
	go dpwss.Start()

	//// Run CSWSServer
	cswss = wss.NewCSWSServer("candlesticks")
	log.Printf("======= CSWSServer[%s] is ready...", cswss.GetName())
	go cswss.Start()

	// Hang thread Main.
	<-finish
}
