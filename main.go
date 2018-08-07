package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ntc-gwss/conf"
	"os"
	"strings"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//// Declare Global
var mapSymbol = make(map[string]string)
var typeTime = map[string]string{
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

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func startWSServer() {
	c := config.GetConfig()

	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	address := c.GetString("ntc.wss.host") + ":" + c.GetString("ntc.wss.port")
	//err := http.ListenAndServe(*addr, nil)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initMapSymbol() {
	c := config.GetConfig()
	listpair := c.GetString("market.listpair")
	log.Printf("=========== listpair: %s", listpair)

	//var listpair = "ETH_BTC;KNOW_BTC;KNOW_ETH"
	var arrSymbol = strings.Split(listpair, ";")
	// log.Printf("arrSymbol: ", arrSymbol)
	for i := range arrSymbol {
		symbol := arrSymbol[i]
		mapSymbol[symbol] = symbol
	}
	log.Printf("=========== mapSymbol: ", mapSymbol)
}

func reloadMapSymbol(listpair string) {
	log.Printf("=========== reloadMapSymbol.listpair: %s", listpair)
	if listpair != "" {
		var arrSymbol = strings.Split(listpair, ";")
		// log.Printf("arrSymbol: ", arrSymbol)
		for i := range arrSymbol {
			symbol := arrSymbol[i]
			//// If not exist, add to mapSymbol
			if _, ok := mapSymbol[symbol]; !ok {
				mapSymbol[symbol] = symbol
			}
		}
		log.Printf("=========== reloadMapSymbol.mapSymbol: ", mapSymbol)
	}
}

func main() {
	// TCF{
	// 	Try: func() {
	// 		fmt.Println("I tried")
	// 		Throw("Oh,...sh...")
	// 	},
	// 	Catch: func(e Exception) {
	// 		fmt.Printf("Caught %v\n", e)
	// 	},
	// 	Finally: func() {
	// 		fmt.Println("Finally...")
	// 	},
	// }.Do()

	environment := flag.String("e", "development", "run project with mode [-e development | test | production]")
	flag.Usage = func() {
		fmt.Println("Usage: [appname] -e development|production")
		os.Exit(1)
	}
	flag.Parse()

	log.Printf("============== environment: %s", *environment)
	config.Init(*environment)

	//// initMapSymbol
	initMapSymbol()

	//// test reloadMapSymbol
	// reloadMapSymbol("ETH_BTC;KNOW_BTC;KNOW_ETH")
	// log.Printf("=========== mapSymbol 2: ", mapSymbol)

	//// Start WSServer.
	//startWSServer()

	//// New WSServer
	// wss := newWSServer("ntc")
	// log.Printf("======= WSServer[%s] is running...", wss.name)

	//// New DPWSServer
	dpwss := newDPWSServer("depthprice")
	log.Printf("======= DPWSServer[%s] is ready...", dpwss.name)
	dpwss.start()
}

// func main() {
// 	environment := flag.String("e", "development", "run project with mode [-e development|production]")
// 	flag.Usage = func() {
// 		fmt.Println("Usage: server -e development|production")
// 		os.Exit(1)
// 	}
// 	flag.Parse()

// 	log.Printf("============== environment: %s", *environment)
// 	config.Init(*environment)

// 	c := config.GetConfig()
// 	// fmt.Println("=======Host: %s", c.GetString("ntc.wss.host"))
// 	// fmt.Println("=======Post: %s", c.GetString("ntc.wss.port"))
// 	log.Printf("=======Host: %s", c.GetString("ntc.wss.host"))
// 	log.Printf("=======Post: %s", c.GetString("ntc.wss.port"))
// }
