package wss

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

func (wss *DPWSServer) GetName() string {
	return wss.name
}

func (wss *DPWSServer) GetHub() *HubLevel1 {
	return wss.hub
}

func NewDPWSServer(name string) *DPWSServer {
	hub := newHubLevel1()
	go hub.run()
	return &DPWSServer{name: name, hub: hub}
}

func (wss *DPWSServer) Start() {
	c := config.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/ws/v1/dp/{symbol}", func(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("WSServer is running on: %s", address)
	log.Printf("======= DPWSServer[%s] is running...", wss.name)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
