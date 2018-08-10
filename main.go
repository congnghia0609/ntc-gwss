package main

import (
	"flag"
	"fmt"
	"log"
	"ntc-gwss/conf"
	"ntc-gwss/server"
	"ntc-gwss/wss"
	"os"

	"github.com/natefinch/lumberjack"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

//// Declare Global
var dpwss *wss.DPWSServer
var htwss *wss.HTWSServer
var cswss *wss.CSWSServer
var tkwss *wss.TKWSServer

// https://github.com/natefinch/lumberjack
func initLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/data/log/ntc-gwss/ntc-gwss.log",
		MaxSize:    10,   // megabytes. Defaults to 100 MB.
		MaxBackups: 3,    // maximum number of old log files to retain.
		MaxAge:     28,   // maximum number of days to retain old log files
		Compress:   true, // disabled by default
	})
}

func main() {
	finish := make(chan bool)

	//// init Configuration
	environment := flag.String("e", "development", "run project with mode [-e development | test | production]")
	flag.Usage = func() {
		fmt.Println("Usage: [appname] -e development | test | production")
		os.Exit(1)
	}
	flag.Parse()
	log.Printf("============== environment: %s", *environment)
	conf.Init(*environment)

	//// init Logger
	if "development" != *environment {
		initLogger()
	}

	//// initMapSymbol
	wss.InitMapSymbol()

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

	// StartWebServer
	go server.StartWebServer("webserver")

	// Hang thread Main.
	<-finish
	log.Println("################# End Main #################")
}
