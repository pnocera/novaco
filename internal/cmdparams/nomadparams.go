package cmdparams

import (
	"runtime"

	"github.com/pnocera/novaco/internal/utils"
)

type NomadConfigParams struct {
	LogLevel        string
	DataDir         string
	BootstrapExpect int
	AdvertiseAddr   string
	BindAddr        string
}

func GetNomadProgramParams(assets string) (*ProgramParams, error) {

	nomaddir := utils.Join(assets, "nomad")

	configtemplate := utils.Join(nomaddir, "templates/nomad.server.hcl")
	configoutput := utils.Join(nomaddir, "config/nomad.server.hcl")

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	nomadconfigparams := NomadConfigParams{
		LogLevel:        "DEBUG",
		DataDir:         utils.Join(nomaddir, "data"),
		BootstrapExpect: 1,
		AdvertiseAddr:   ip.String(),
		BindAddr:        "0.0.0.0",
	}

	err = utils.Render(configtemplate, configoutput, nomadconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(nomaddir, "nomad_"+runtime.GOARCH+".exe")

	return &ProgramParams{
		DirPath:     nomaddir,
		ExeFullname: exefile,
		AdditionalParams: []string{
			"agent",
			"-config=" + configoutput,
		},
	}, nil
}
