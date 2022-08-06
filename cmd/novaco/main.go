package main

import (
	"log"

	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/program"
)

var logger service.Logger

func main() {
	serviceConfig := &service.Config{
		Name:        "gci-nomad",
		DisplayName: "GCI Nomad Server",
		Description: "A Hashicorp Nomad Primary server",
	}

	prg := &program.Program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.Service = s

	logger, err = s.Logger(nil)

	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
