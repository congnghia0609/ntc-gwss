/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"github.com/congnghia0609/ntc-gwss/conf"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// CSWSServer class
type CSWSServer struct {
	name string
	hub  *HubLevel2
}

var mapInstanceCS = make(map[string]*CSWSServer)

// GetInstanceCS get instance CS
func GetInstanceCS(name string) *CSWSServer {
	return mapInstanceCS[name]
}

// GetName get name
func (wss *CSWSServer) GetName() string {
	return wss.name
}

// GetHub get hub
func (wss *CSWSServer) GetHub() *HubLevel2 {
	return wss.hub
}

// NewCSWSServer new CSWSServer
func NewCSWSServer(name string) *CSWSServer {
	hub := newHubLevel2()
	go hub.run()
	instance := &CSWSServer{name: name, hub: hub}
	mapInstanceCS[name] = instance
	return instance
}

// Start CSWSServer
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
	log.Printf("======= CSWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("CSWSServer ListenAndServe: ", err)
	}
}
