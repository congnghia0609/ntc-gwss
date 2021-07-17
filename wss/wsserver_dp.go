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

// DPWSServer class
type DPWSServer struct {
	name string
	hub  *HubLevel1
}

var mapInstanceDP = make(map[string]*DPWSServer)

// GetInstanceDP get instance DP
func GetInstanceDP(name string) *DPWSServer {
	return mapInstanceDP[name]
}

// GetName get name
func (wss *DPWSServer) GetName() string {
	return wss.name
}

// GetHub get hub
func (wss *DPWSServer) GetHub() *HubLevel1 {
	return wss.hub
}

// NewDPWSServer new DPWSServer
func NewDPWSServer(name string) *DPWSServer {
	hub := newHubLevel1()
	go hub.run()
	instance := &DPWSServer{name: name, hub: hub}
	mapInstanceDP[name] = instance
	return instance
}

// Start DPWSServer
func (wss *DPWSServer) Start() {
	c := conf.GetConfig()

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
	log.Printf("======= DPWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("DPWSServer ListenAndServe: ", err)
	}
}
