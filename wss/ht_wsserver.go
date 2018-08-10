package wss

import (
	"log"
	"net/http"
	"ntc-gwss/conf"
	"strings"

	"github.com/gorilla/mux"
)

// HTWSServer class
type HTWSServer struct {
	name string
	hub  *HubLevel1
}

var mapInstanceHT = make(map[string]*HTWSServer)

func GetInstanceHT(name string) *HTWSServer {
	return mapInstanceHT[name]
}

func (wss *HTWSServer) GetName() string {
	return wss.name
}

func (wss *HTWSServer) GetHub() *HubLevel1 {
	return wss.hub
}

func NewHTWSServer(name string) *HTWSServer {
	hub := newHubLevel1()
	go hub.run()
	instance := &HTWSServer{name: name, hub: hub}
	mapInstanceHT[name] = instance
	return instance
}

func (wss *HTWSServer) Start() {
	c := conf.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/ws/v1/ht/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		if len(vars["symbol"]) > 0 {
			symbol := vars["symbol"]
			symbol = strings.ToUpper(symbol)
			log.Printf("=======symbol: %s", symbol)
			if _, ok := MapSymbol[symbol]; ok {
				serveWsLevel1(symbol, wss.hub, w, r)
			}
		}
	})
	httpsm.Handle("/", rt)

	address := c.GetString(wss.name+".wss.host") + ":" + c.GetString(wss.name+".wss.port")
	// log.Printf("WSServer is running on: %s", address)
	log.Printf("======= HTWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
