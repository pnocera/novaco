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
	ConsulPort     string
	GitPort        string
	GitHost        string
	ApiPort        string
	ApiHost        string
}

// GetConsulProgramParams returns the program params for consul
func GetConsulProgramParams() (*ProgramParams, error) {

	ip := sets.IP()

	consulconfigparams := ConsulConfigParams{
		LogLevel:       sets.UppercaseLogLevel(),
		DataDir:        utils.DataPath("consul"),
		DataCenter:     "dc1",
		ClientAddr:     ip,
		UiEnabled:      true,
		BindAddr:       ip,
		Server:         true,
		Bootstrap:      true,
		ConnectEnabled: true,
		AddressesHttp:  ip + " 127.0.0.1",
		ConsulPort:     sets.ConsulPort,
		GitPort:        sets.GitPort,
		GitHost:        ip,
		ApiPort:        sets.APIPort,
		ApiHost:        ip,
	}

	configtemplate := utils.TemplatePath("consul.hcl")
	configdir := utils.ConfigPath("consul")
	configoutput := utils.Join(configdir, "consul.hcl")
	err := utils.Render(configtemplate, configoutput, consulconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(utils.BinPath("consul"), "consul_"+runtime.GOARCH+".exe")

	additionalparams := []string{
		"agent",
		"-config-file=" + configdir,
	}

	return &ProgramParams{
		ID:               "consul",
		DirPath:          filepath.Dir(exefile),
		ExeFullname:      exefile,
		AdditionalParams: additionalparams,
	}, nil
}
