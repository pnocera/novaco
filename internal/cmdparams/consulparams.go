package cmdparams

import (
	"path/filepath"
	"runtime"

	"github.com/pnocera/novaco/internal/utils"
)

type ConsulConfigParams struct {
	LogLevel       string
	DataDir        string
	DataCenter     string
	ClientAddr     string
	UiEnabled      bool
	Server         bool
	BindAddr       string
	Bootstrap      bool
	ConnectEnabled bool
	AddressesHttp  string
}

// GetConsulProgramParams returns the program params for consul
// assets: the path to the assets directory
// runtype is the type of run (primary, server or client)
func GetConsulProgramParams(assets string, runtype string) (*ProgramParams, error) {

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	consulconfigparams := ConsulConfigParams{
		LogLevel:       "DEBUG",
		DataDir:        utils.Join(assets, "data/consul"),
		DataCenter:     "dc1",
		ClientAddr:     ip.String(),
		UiEnabled:      true,
		BindAddr:       ip.String(),
		Server:         true,
		Bootstrap:      true,
		ConnectEnabled: true,
		AddressesHttp:  ip.String() + " 127.0.0.1",
	}

	configtemplate := utils.Join(assets, "templates/"+runtype+"/consul.hcl")
	configoutput := utils.Join(assets, "config/consul/auto/consul."+runtype+".hcl")
	err = utils.Render(configtemplate, configoutput, consulconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(assets, "bin/consul/consul_"+runtime.GOARCH+".exe")

	additionalparams := []string{
		"agent",
		"-config-file=" + configoutput,
		"-config-file=" + utils.Join(assets, "config/consul/custom"),
	}

	if runtype == "dev" {
		additionalparams = append(additionalparams, "-dev")
	}

	return &ProgramParams{
		DirPath:          filepath.Dir(exefile),
		ExeFullname:      exefile,
		AdditionalParams: additionalparams,
		LogFile:          utils.Join(assets, "logs/consul/consul."+runtype+".log"),
	}, nil
}
