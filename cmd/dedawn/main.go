package main

import (
	"dedawn/mask"
	"flag"
	"gioui.org/app"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	server := "http://localhost:8080"
	logpath := "dedawn.log"
	flag.StringVar(&server, "server", "http://localhost:8080", "the server address")
	flag.StringVar(&logpath, "logpath", "dedawn.log", "the log path")
	flag.Parse()
	logfile, err := os.Create(logpath)
	if err != nil {
		log.Errorf("create log file %v", err)
		return
	}
	log.SetOutput(logfile)

	mask.Run(server)

	app.Main()
}
