package service

import (
	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/newprogram"
	"github.com/pnocera/novaco/internal/settings"
)

var sets = settings.GetSettings()

func StartNew() error {

	svcConfig := &service.Config{
		Name:        "gci-nomad",
		DisplayName: "GCI Nomad Server ",
		Description: "A Hashicorp Nomad server",
	}

	prg := newprogram.NewProgram()
	s, err := service.New(prg, svcConfig)
	if err != nil {
		sets.Logger.Error("Error creating service: %v", err)
		return err
	}

	err = s.Run()

	if err != nil {
		sets.Logger.Error("Error running service: %v", err)
	}
	sets.Logger.Info("Service stopped")
	return err
}
