package main

import (
	"log"
	"net/http"
	"ntc-gwss/conf"

	"github.com/gorilla/mux"
)

// WSServer class
type WSServer struct {
	name string
	hub  *Hub
}

func newWSServer(name string) *WSServer {
	c := config.GetConfig()

	hub := newHub()
	go hub.run()
	rt := mux.NewRouter()
	// http.HandleFunc("/", serveHome)
	// http.HandleFunc("/ws/v1/dp/:symbol", func(w http.ResponseWriter, r *http.Request) {
	// 	pathAddr := r.RemoteAddr
	// 	log.Printf("=======pathAddr: %s", pathAddr)
	// 	pathURI := r.RequestURI
	// 	log.Printf("=======pathURI: %s", pathURI)

	// 	serveWs(hub, w, r)
	// })
	rt.HandleFunc("/", serveHome)
	rt.HandleFunc("/ws/v1/dp/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		log.Printf("=======symbol: %s", symbol)

		serveWs(hub, w, r)
	})
	http.Handle("/", rt)

	address := c.GetString(name+".wss.host") + ":" + c.GetString(name+".wss.port")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return &WSServer{name: name, hub: hub}
}
