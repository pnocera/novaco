package service

import (
	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/klogger"
	"github.com/pnocera/novaco/internal/newprogram"
)

var klog = klogger.NewKLogger("service")

func StartNew(runtype string) error {

	svcConfig := &service.Config{
		Name:        "gci-nomad-" + runtype,
		DisplayName: "GCI Nomad Server " + runtype,
		Description: "A Hashicorp Nomad server",
	}

	prg := newprogram.NewProgram(runtype)
	s, err := service.New(prg, svcConfig)
	if err != nil {
		klog.Error("Error creating service: %v", err)
		return err
	}

	err = s.Run()

	if err != nil {
		klog.Error("Error running service: %v", err)
	}
	klog.Info("Service stopped")
	return err
}
