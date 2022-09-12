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

func GetNomadProgramParams() (*ProgramParams, error) {

	nomadconfigparams := NomadConfigParams{
		LogLevel:        sets.UppercaseLogLevel(),
		DataDir:         utils.DataPath("nomad"),
		BootstrapExpect: 1,
		AdvertiseAddr:   sets.IP(),
		BindAddr:        "0.0.0.0",
	}

	configtemplate := utils.TemplatePath("nomad.hcl")
	configoutput := utils.Join(utils.ConfigPath("nomad"), "nomad.hcl")

	err := utils.Render(configtemplate, configoutput, nomadconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(utils.BinPath("nomad"), "nomad_"+runtime.GOARCH+".exe") //  "nomad.exe") //

	additionalparams := []string{
		"agent",
		"-config=" + utils.ConfigPath("nomad"),
	}

	return &ProgramParams{
		ID:               "nomad",
		DirPath:          filepath.Dir(exefile),
		ExeFullname:      exefile,
		AdditionalParams: additionalparams,
	}, nil
}
