package cmdparams

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/utils"
)

var sets = settings.GetSettings()

func RenderGitConfigIfNotExist() error {

	datapath := utils.DataPath("gitea")
	gitconfigparams := config.GitConfig{
		LogLevel:     sets.FirstUppercaseLogLevel(),
		LogMode:      "console",
		HostIP:       utils.IP(),
		Port:         sets.GitPort,
		GitPath:      utils.Join(utils.BinPath("gitea"), "git.exe"),
		DatabasePath: utils.Join(datapath, "gitea.db"),
		Domain:       utils.IP(),
		RunMode:      "prod",
		RunUser:      "COMPUTERNAME$",
		RepoPath:     utils.Join(datapath, "repo"),
		LfsPath:      utils.Join(datapath, "lfs"),
	}

	configtemplate := utils.TemplatePath("gitea.ini")
	configoutput := utils.Join(utils.ConfigPath("gitea"), "gitea.ini")

	if _, err := os.Stat(configoutput); errors.Is(err, os.ErrNotExist) {
		err = utils.Render(configtemplate, configoutput, gitconfigparams)
		return err
	}

	return nil
}

func GetGitParams() (*ProgramParams, error) {

	exefile := utils.Join(utils.BinPath("gitea"), "gitea.exe")

	return &ProgramParams{
		ID:          "gitea",
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"web",
			"-c",
			utils.Join(utils.ConfigPath("gitea"), "gitea.ini"),
			"-w",
			utils.DataPath("gitea"),
			"--verbose",
		},
	}, nil
}
