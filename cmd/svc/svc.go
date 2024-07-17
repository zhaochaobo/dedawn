package main

import (
	"dedawn/mask"
	"flag"
	"fmt"
	"gioui.org/app"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

func pro() {

	ps, err := process.Processes()
	if err != nil {
		log.Errorf("get processes %v", err)
		return
	}

	for _, p := range ps {
		n, err := p.Name()
		if err != nil {
			log.Errorf("get process name %d, %v", p.Pid, err)
			return
		}
		log.Infof("process: %v", n)
	}

}

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	server := "http://localhost:8080"
	flag.StringVar(&server, "server", "http://localhost:8080", "the server address")
	flag.Parse()
	fmt.Println(server)

	mask.Run(server)
	log.Printf("mask run with server %s", server)

	app.Main()
	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "DedawnService",
		DisplayName: "Dedawn Servcie",
		Description: "This is dedawn servcie.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
