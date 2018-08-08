package wss

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
	http.ServeFile(w, r, "view/home.html")
}

func NewWSServer(name string) *WSServer {
	c := config.GetConfig()

	hub := newHub()
	go hub.run()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/", serveHome)
	rt.HandleFunc("/ws/v1/dp/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		log.Printf("=======symbol: %s", symbol)

		serveWs(hub, w, r)
	})
	httpsm.Handle("/", rt)

	address := c.GetString(name+".wss.host") + ":" + c.GetString(name+".wss.port")
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return &WSServer{name: name, hub: hub}
}
