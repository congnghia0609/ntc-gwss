/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"log"
	"net/http"
	"strings"
	"ntc-gwss/conf"

	"github.com/gorilla/mux"
)

// CRWSServer class
type CRWSServer struct {
	name string
	hub  *HubCR
}

var mapInstanceCR = make(map[string]*CRWSServer)

func GetInstanceCR(name string) *CRWSServer {
	return mapInstanceCR[name]
}

func (wss *CRWSServer) GetName() string {
	return wss.name
}

func (wss *CRWSServer) GetHub() *HubCR {
	return wss.hub
}

func NewCRWSServer(name string) *CRWSServer {
	hub := newHubCR()
	go hub.run()
	instance := &CRWSServer{name: name, hub: hub}
	mapInstanceCR[name] = instance
	return instance
}

func (wss *CRWSServer) Start() {
	c := conf.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/ws/v1/cr/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		if len(vars["symbol"]) > 0 {
			symbol := vars["symbol"]
			symbol = strings.ToUpper(symbol)
			log.Printf("=======symbol: %s", symbol)
			if _, ok := MapSymbol[symbol]; ok {
				serveWsCR(symbol, wss.hub, w, r)
			}
		}
	})
	rt.HandleFunc("/ws/v2/cr/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		vars := mux.Vars(r)
		if len(vars["symbol"]) > 0 {
			symbol := vars["symbol"]
			symbol = strings.ToUpper(symbol)
			log.Printf("=======symbol: %s", symbol)
			if _, ok := MapSymbol[symbol]; ok {
				serveWsCR(symbol, wss.hub, w, r)
			}
		}
	})
	httpsm.Handle("/", rt)

	address := c.GetString(wss.name+".wss.host") + ":" + c.GetString(wss.name+".wss.port")
	log.Printf("======= CRWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("CRWSServer ListenAndServe: ", err)
	}
}
