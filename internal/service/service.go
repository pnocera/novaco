package service

import (
	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/newprogram"
	"github.com/pnocera/novaco/internal/utils"
)

var logger = utils.NewKLogger("service")

func StartNew() error {

	svcConfig := &service.Config{
		Name:        "gci-nomad",
		DisplayName: "GCI Nomad Server ",
		Description: "A Hashicorp Nomad server",
	}

	prg := newprogram.NewProgram()
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error("Error creating service: %v", err)
		return err
	}

	err = s.Run()

	if err != nil {
		logger.Error("Error running service: %v", err)
	}
	logger.Info("Service stopped")
	return err
}
