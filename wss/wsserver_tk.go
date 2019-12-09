/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wss

import (
	"log"
	"net/http"
	"ntc-gwss/conf"

	"github.com/gorilla/mux"
)

// TKWSServer class
type TKWSServer struct {
	name string
	hub  *Hub
}

var mapInstanceTK = make(map[string]*TKWSServer)
var TKDataCache string
var TKDataCacheCR string

func GetInstanceTK(name string) *TKWSServer {
	return mapInstanceTK[name]
}

func (wss *TKWSServer) GetName() string {
	return wss.name
}

func (wss *TKWSServer) GetHub() *Hub {
	return wss.hub
}

func NewTKWSServer(name string) *TKWSServer {
	hub := newHub()
	go hub.run()
	instance := &TKWSServer{name: name, hub: hub}
	mapInstanceTK[name] = instance
	return instance
}

func (wss *TKWSServer) Start() {
	c := conf.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()
	rt.HandleFunc("/ws/v1/tk", func(w http.ResponseWriter, r *http.Request) {
		pathURI := r.RequestURI
		log.Printf("=======pathURI: %s", pathURI)
		serveWs(wss.hub, w, r)
	})
	httpsm.Handle("/", rt)

	address := c.GetString(wss.name+".wss.host") + ":" + c.GetString(wss.name+".wss.port")
	log.Printf("======= TKWSServer[%s] is running on host: %s", wss.name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("TKWSServer ListenAndServe: ", err)
	}
}
