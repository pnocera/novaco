package cmdparams

import (
	"path/filepath"

	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/utils"
)

func GetGitParams(assets string, runtype string) (*ProgramParams, error) {

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	gitconfigparams := config.Config{
		LogLevel:            "DEBUG",
		ShutdownTimeout:     "15s",
		Hostname:            ip.String(),
		Port:                8888,
		SSLEnabled:          false,
		AuthEnabled:         false,
		PasswdFilePath:      utils.Join(assets, "passwd"),
		RestrictReceivePack: false,
		RestrictUploadPack:  false,

		GitBinPath: utils.Join(assets, "bin/git"),
		ReposPath:  utils.Join(assets, "data/git"),
	}

	configtemplate := utils.Join(assets, "templates/git."+runtype+".hcl")
	configoutput := utils.Join(assets, "config/git/auto/git."+runtype+".hcl")

	err = utils.Render(configtemplate, configoutput, gitconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(assets, "bin/git/git-server.exe")

	return &ProgramParams{
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"-config=" + configoutput,
			"-config=" + utils.Join(assets, "config/git/custom"),
		},
		LogFile: utils.Join(assets, "logs/git/git."+runtype+".log"),
	}, nil
}
