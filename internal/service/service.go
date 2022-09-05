package service

import (
	"log"

	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/newprogram"
)

func StartNew(runtype string) {

	svcConfig := &service.Config{
		Name:        "gci-nomad-" + runtype,
		DisplayName: "GCI Nomad Server " + runtype,
		Description: "A Hashicorp Nomad server",
	}

	prg := newprogram.NewProgram(runtype)
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Printf("Error creating new service %v", err)
	}

	err = s.Run()
	log.Println("Service stopped")
	if err != nil {
		log.Printf("Error running service %v", err)
	}
}
