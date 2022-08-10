package cmdparams

import (
	"path/filepath"
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

func GetNomadProgramParams(assets string, runtype string) (*ProgramParams, error) {

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	nomadconfigparams := NomadConfigParams{
		LogLevel:        "DEBUG",
		DataDir:         utils.Join(assets, "data/nomad"),
		BootstrapExpect: 1,
		AdvertiseAddr:   ip.String(),
		BindAddr:        "0.0.0.0",
	}

	configtemplate := utils.Join(assets, "templates/nomad."+runtype+".hcl")
	configoutput := utils.Join(assets, "config/nomad/auto/nomad."+runtype+".hcl")

	err = utils.Render(configtemplate, configoutput, nomadconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(assets, "bin/nomad/nomad_"+runtime.GOARCH+".exe")

	return &ProgramParams{
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"agent",
			"-config=" + configoutput,
			"-config=" + utils.Join(assets, "config/nomad/custom"),
		},
		LogFile: utils.Join(assets, "logs/nomad/nomad."+runtype+".log"),
	}, nil
}
