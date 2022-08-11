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

	gitconfigparams := config.GitConfig{
		LogLevel:   "DEBUG",
		Hostname:   ip.String(),
		Port:       8888,
		KeyDir:     "",
		Dir:        utils.Join(assets, "data/git"),
		GitPath:    utils.Join(assets, "bin/git"),
		GitUser:    "git",
		AutoCreate: true,
		AutoHooks:  false,
		Hooks:      &config.HookScripts{},
		Auth:       false,
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
