package main

import (
	"flag"
	"fmt"
	"log"
	"ntc-gwss/conf"
	"ntc-gwss/wss"
	"os"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//// Declare Global
var dpwss *wss.DPWSServer
var cswss *wss.CSWSServer

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
	wss.InitMapSymbol()

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
