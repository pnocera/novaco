package cmdparams

import (
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

func GetConsulProgramParams(assets string) (*ProgramParams, error) {

	consuldir := utils.Join(assets, "consul")

	configtemplate := utils.Join(consuldir, "templates/consul.server.hcl")
	configoutput := utils.Join(consuldir, "config/consul.server.hcl")

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	consulconfigparams := ConsulConfigParams{
		LogLevel:       "DEBUG",
		DataDir:        utils.Join(consuldir, "data"),
		DataCenter:     "dc1",
		ClientAddr:     ip.String(),
		UiEnabled:      true,
		BindAddr:       ip.String(),
		Server:         true,
		Bootstrap:      true,
		ConnectEnabled: true,
		AddressesHttp:  ip.String() + " 127.0.0.1",
	}

	err = utils.Render(configtemplate, configoutput, consulconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(consuldir, "consul_"+runtime.GOARCH+".exe")

	return &ProgramParams{
		DirPath:     consuldir,
		ExeFullname: exefile,
		AdditionalParams: []string{
			"agent",
			"-config-file=" + configoutput,
		},
	}, nil
}
