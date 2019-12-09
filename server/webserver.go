/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ntc-gwss/conf"
	"ntc-gwss/wss"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

func printJson(w http.ResponseWriter, r *http.Request, data string) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Response to client.
	io.WriteString(w, data)
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
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

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// http://localhost:15901/tmws/v1/status
func statusHandle(w http.ResponseWriter, r *http.Request) {
	pathURI := r.RequestURI
	log.Printf("=======pathURI: %s", pathURI)
	// vars := mux.Vars(r)
	// symbol := vars["symbol"]
	// log.Printf("=======symbol: %s", symbol)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mapData := make(map[string]string)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memmory // For info on each, see: https://golang.org/pkg/runtime/#MemStats
	mapData["A_Alloc"] = fmt.Sprint(bToMb(m.Alloc)) + " MB"
	// mapData["A_TotalAlloc"] = fmt.Sprint(bToMb(m.TotalAlloc)) + " MB"
	mapData["A_Sys"] = fmt.Sprint(bToMb(m.Sys)) + " MB"
	mapData["A_StackSys"] = fmt.Sprint(bToMb(m.StackSys)) + " MB"

	// MapSymbol
	// mapMarketName, _ := json.Marshal(wss.MapSymbol)
	mapData["B_mapMarketName"] = fmt.Sprint(wss.MapSymbol)

	mapTT := make(map[string]int)
	// Count client
	for sk := range wss.MapSymbol {
		// log.Printf("SKey=%s, SValue=%s", sk, sv)
		// DPWSServer
		dpws := wss.GetInstanceDP(wss.NameDPWSS)
		if dpws != nil {
			keyDP := fmt.Sprintf("DPWSServer.clients[%s].size", sk)
			valueDP := dpws.GetHub().GetSizeClientLevel1(sk)
			mapData[keyDP] = fmt.Sprint(valueDP)
		}
		// HTWSServer
		htws := wss.GetInstanceHT(wss.NameHTWSS)
		if htws != nil {
			keyHT := fmt.Sprintf("HTWSServer.clients[%s].size", sk)
			valueHT := htws.GetHub().GetSizeClientLevel1(sk)
			mapData[keyHT] = fmt.Sprint(valueHT)
		}
		// CSWSServer
		csws := wss.GetInstanceCS(wss.NameCSWSS)
		if csws != nil {
			for tk := range wss.TypeTime {
				key := sk + "_" + tk
				// keyCS := fmt.Sprintf("CSWSServer.clients[%s].size", key)
				valueCS := csws.GetHub().GetSizeClientLevel2(key)
				// mapData[keyCS] = fmt.Sprint(valueCS)
				if _, ok := mapTT[tk]; ok {
					ctt := mapTT[tk] + valueCS
					mapTT[tk] = ctt
				} else {
					mapTT[tk] = valueCS
				}
			}
		}
	}
	// CSWSServer
	mapData["StockWSSServer.mapTT"] = fmt.Sprint(mapTT)

	// TKWSServer
	tkwss := wss.GetInstanceTK(wss.NameTKWSS)
	if tkwss != nil {
		keyTK := "TKWSServer.clients.size"
		valueTK := tkwss.GetHub().GetSizeClient()
		mapData[keyTK] = fmt.Sprint(valueTK)
	}

	// CRWSServer
	crwss := wss.GetInstanceCR(wss.NameCRWSS)
	if crwss != nil {
		keyCR := "ZCRWSServer.clients.size"
		valueCR := crwss.GetHub().GetSizeClientCR()
		mapData[keyCR] = fmt.Sprint(valueCR)
	}

	// timestamp
	now := time.Now()
	timestamp := now.UnixNano() / 1000000
	mapData["timestamp"] = fmt.Sprint(timestamp)
	timeserver := fmt.Sprint(now.Format(time.RFC3339))
	mapData["time_server"] = fmt.Sprint(timeserver)

	// log.Printf("mapData: ", mapData)
	data, _ := json.Marshal(mapData)
	// Response.
	if len(data) > 0 {
		printJson(w, r, string(data))
	} else {
		printJson(w, r, "{}")
	}
}

// http://localhost:15901/tmws/v1/as?list_symbol=BTC_USDT%3BETH_USDT%3BKNOW_USDT%3BGTO_USDT
func statusAddSymbol(w http.ResponseWriter, r *http.Request) {
	mapData := make(map[string]string)
	mapData["err"] = "-1"
	mapData["msg"] = "Execute fail."

	pathURI := r.RequestURI
	log.Printf("=======pathURI: %s", pathURI)
	// vars := mux.Vars(r)
	// list_symbol := vars["list_symbol"]
	// log.Printf("=======list_symbol: %s", list_symbol)

	listsymbol := r.FormValue("list_symbol")
	log.Printf("=======listsymbol: %s", listsymbol)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(listsymbol) > 0 {
		wss.ReloadMapSymbol(listsymbol)
		mapData["err"] = "0"
		mapData["msg"] = "Add symbol successfully."
	}

	// log.Printf("mapData: ", mapData)
	data, _ := json.Marshal(mapData)
	// Response.
	if len(data) > 0 {
		printJson(w, r, string(data))
	} else {
		printJson(w, r, "{}")
	}
}

func StartWebServer(name string) {
	c := conf.GetConfig()

	// NewServeMux
	httpsm := http.NewServeMux()

	// Setup Handlers.
	rt := mux.NewRouter()

	// static resources
	// This will serve files under http://localhost:15901/static/<filename>
	rt.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	// Mapping Handlers
	// rt.HandleFunc("/", homeHandle)
	rt.HandleFunc("/tmws/v1/status", statusHandle)
	rt.HandleFunc("/tmws/v1/as", statusAddSymbol)
	httpsm.Handle("/", rt)

	address := c.GetString(name+".host") + ":" + c.GetString(name+".port")
	log.Printf("======= WebServer[%s] is running on host: %s", name, address)
	err := http.ListenAndServe(address, httpsm)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
