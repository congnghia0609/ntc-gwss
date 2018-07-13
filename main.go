package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"ntc-gwss/conf"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

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

func main() {
	environment := flag.String("e", "development", "run project with mode [-e development|production]")
	flag.Usage = func() {
		fmt.Println("Usage: [appname] -e development|production")
		os.Exit(1)
	}
	flag.Parse()

	log.Printf("============== environment: %s", *environment)
	config.Init(*environment)

	// Start WSServer.
	//startWSServer()

	// New WSServer
	wss := newWSServer("ntc")
	log.Printf("======= WSServer[%s] is running...", wss.name)
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
