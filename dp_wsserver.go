package main

import (
	"log"
	"net/http"
	"ntc-gwss/conf"
	"strings"

	"github.com/gorilla/mux"
)

// DPWSServer class
type DPWSServer struct {
	name string
	hub  *HubLevel1
}

func newDPWSServer(name string) *DPWSServer {
	// c := config.GetConfig()

	hub := newHubLevel1()
	go hub.run()

	// // Setup Handlers.
	// rt := mux.NewRouter()
	// rt.HandleFunc("/", serveHome)
	// rt.HandleFunc("/ws/v1/dp/{symbol}", func(w http.ResponseWriter, r *http.Request) {
	// 	pathURI := r.RequestURI
	// 	log.Printf("=======pathURI: %s", pathURI)
	// 	vars := mux.Vars(r)
	// 	if len(vars["symbol"]) > 0 {
	// 		symbol := vars["symbol"]
	// 		symbol = strings.ToUpper(symbol)
	// 		log.Printf("=======symbol: %s", symbol)
	// 		if _, ok := mapSymbol[symbol]; ok {
	// 			serveWsLevel1(symbol, hub, w, r)
	// 		}
	// 	}
	// })
	// http.Handle("/", rt)

	// address := c.GetString(name+".wss.host") + ":" + c.GetString(name+".wss.port")
	// log.Printf("WSServer is running on: %s", address)
	// err := http.ListenAndServe(address, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	return &DPWSServer{name: name, hub: hub}
}

func (wss DPWSServer) start() {
	c := config.GetConfig()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/", serveHome)
	rt.HandleFunc("/ws/v1/dp/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		if len(vars["symbol"]) > 0 {
			symbol := vars["symbol"]
			symbol = strings.ToUpper(symbol)
			log.Printf("=======symbol: %s", symbol)
			if _, ok := mapSymbol[symbol]; ok {
				serveWsLevel1(symbol, wss.hub, w, r)
			}
		}
	})
	http.Handle("/", rt)

	address := c.GetString(wss.name+".wss.host") + ":" + c.GetString(wss.name+".wss.port")
	log.Printf("WSServer is running on: %s", address)
	log.Printf("======= DPWSServer[%s] is running...", wss.name)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
