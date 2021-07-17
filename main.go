/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package main

import (
	"flag"
	"fmt"
	"github.com/congnghia0609/ntc-gwss/conf"
	"github.com/congnghia0609/ntc-gwss/server"
	"github.com/congnghia0609/ntc-gwss/wss"
	"log"
	"os"
	"os/signal"

	"github.com/natefinch/lumberjack"
)

//// Declare Global
// WSServer
var dpwss *wss.DPWSServer
var htwss *wss.HTWSServer
var cswss *wss.CSWSServer
var tkwss *wss.TKWSServer
var crwss *wss.CRWSServer

// WSClient
//var dpwsc *wsc.NWSClient
//var cswsc *wsc.NWSClient
//var htwsc *wsc.NWSClient
//var tkwsc *wsc.NWSClient
//var crwsc *wsc.NWSClient
//var rswsc *wsc.NWSClient

// https://github.com/natefinch/lumberjack
func initLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/data/log/ntc-gwss/ntc-gwss.log",
		MaxSize:    10,   // 10 megabytes. Defaults to 100 MB.
		MaxBackups: 3,    // maximum number of old log files to retain.
		MaxAge:     28,   // maximum number of days to retain old log files
		Compress:   true, // disabled by default
	})
}

func main() {
	////// -------------------- Init System -------------------- //////
	//// init Configuration
	environment := flag.String("e", "development", "run project with mode [-e development | test | production]")
	flag.Usage = func() {
		fmt.Println("Usage: ./[appname] -e development | test | production")
		os.Exit(1)
	}
	flag.Parse()
	log.Printf("============== environment: %s", *environment)
	conf.Init(*environment)

	//// init Logger
	if "development" != *environment {
		log.Printf("============== LogFile: /data/log/ntc-gwss/ntc-gwss.log")
		initLogger()
	}

	//// initMapSymbol
	wss.InitMapSymbol()

	////// -------------------- Start WSServer -------------------- //////
	//// Run DPWSServer
	dpwss = wss.NewDPWSServer(wss.NameDPWSS)
	log.Printf("======= DPWSServer[%s] is ready...", dpwss.GetName())
	go dpwss.Start()

	//// Run HTWSServer
	htwss = wss.NewHTWSServer(wss.NameHTWSS)
	log.Printf("======= HTWSServer[%s] is ready...", htwss.GetName())
	go htwss.Start()

	//// Run CSWSServer
	cswss = wss.NewCSWSServer(wss.NameCSWSS)
	log.Printf("======= CSWSServer[%s] is ready...", cswss.GetName())
	go cswss.Start()

	//// Run TKWSServer
	tkwss = wss.NewTKWSServer(wss.NameTKWSS)
	log.Printf("======= TKWSServer[%s] is ready...", tkwss.GetName())
	go tkwss.Start()

	//// Run CRWSServer
	crwss = wss.NewCRWSServer(wss.NameCRWSS)
	log.Printf("======= CRWSServer[%s] is ready...", crwss.GetName())
	go crwss.Start()

	////// -------------------- Start WSClient -------------------- //////
	//// // DPWSClient
	//dpwsc = wsc.NewDPWSClient()
	//defer dpwsc.Close()
	//go dpwsc.StartDPWSClient()
	//
	//// // CSWSClient
	//cswsc = wsc.NewCSWSClient()
	//defer cswsc.Close()
	//go cswsc.StartCSWSClient()
	//
	//// // HTWSClient
	//htwsc = wsc.NewHTWSClient()
	//defer htwsc.Close()
	//go htwsc.StartHTWSClient()
	//
	//// // TKWSClient
	//tkwsc = wsc.NewTKWSClient()
	//defer tkwsc.Close()
	//go tkwsc.StartTKWSClient()
	//
	//// // CRWSClient
	//crwsc = wsc.NewCRWSClient()
	//defer crwsc.Close()
	//go crwsc.StartCRWSClient()

	// // // ReloadSymbolWSSClient
	// rswsc = wsc.NewRSWSClient()
	// defer rswsc.Close()
	// go rswsc.StartRSWSClient()

	////// -------------------- Start WebServer -------------------- //////
	// StartWebServer
	go server.StartWebServer("webserver")

	// Hang thread Main.
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("################# End Main #################")
}
