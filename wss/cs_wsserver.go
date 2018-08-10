package wss

import (
	"log"
	"net/http"
	"ntc-gwss/conf"
	"strings"

	"github.com/gorilla/mux"
)

// CSWSServer class
type CSWSServer struct {
	name string
	hub  *HubLevel2
}

var mapInstanceCS = make(map[string]*CSWSServer)

func GetInstanceCS(name string) *CSWSServer {
	return mapInstanceCS[name]
}

func (wss *CSWSServer) GetName() string {
	return wss.name
}

func (wss *CSWSServer) GetHub() *HubLevel2 {
	return wss.hub
}

func NewCSWSServer(name string) *CSWSServer {
	hub := newHubLevel2()
	go hub.run()
	instance := &CSWSServer{name: name, hub: hub}
	mapInstanceCS[name] = instance
	return instance
}

func (wss *CSWSServer) Start() {
	c := conf.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/ws/v1/cs/{STT}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		if len(vars["STT"]) > 0 {
			stt := vars["STT"]
			arrSTT := strings.Split(stt, "@")
			if len(arrSTT) >= 2 {
				symbol := strings.ToUpper(arrSTT[0])
				tt := arrSTT[1]
				log.Printf("=======symbol: %s, typeTime: %s", symbol, tt)
				_, ok1 := MapSymbol[symbol]
				_, ok2 := TypeTime[tt]
				if ok1 && ok2 {
					serveWsLevel2(symbol, tt, wss.hub, w, r)
				}
			}
		}
	})
	httpsm.Handle("/", rt)

	address := c.GetString(wss.name+".wss.host") + ":" + c.GetString(wss.name+".wss.port")
	// log.Printf("WSServer is running on: %s", address)
	log.Printf("======= CSWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
