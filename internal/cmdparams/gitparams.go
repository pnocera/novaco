package cmdparams

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/utils"
)

func GetGitIniFilePath(runtype string) (string, error) {
	assets, err := utils.Assets()
	if err != nil {
		return "", err
	}
	return utils.Join(assets, "config/git/app.ini"), nil
}

func RenderGitConfigIfNotExist(runtype string) error {
	assets, err := utils.Assets()
	if err != nil {
		return err
	}

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return err
	}

	runmode := "dev"

	if runtype == "primary" {
		runmode = "prod"
	}

	gitconfigparams := config.GitConfig{
		LogLevel:     "debug",
		LogPath:      utils.Join(assets, "logs/git/git."+runtype+".log"),
		LogMode:      "console",
		HostIP:       ip.String(),
		Port:         8888,
		GitPath:      utils.Join(assets, "bin/git/git.exe"),
		DatabasePath: utils.Join(assets, "data/git/git."+runtype+".db"),
		Domain:       ip.String(),
		RunMode:      runmode,
		RunUser:      "COMPUTERNAME$",
		RepoPath:     utils.Join(assets, "data/git/repo"),
		LfsPath:      utils.Join(assets, "data/git/lfs"),
	}

	configtemplate := utils.Join(assets, "templates/"+runtype+"/git.ini")
	configoutput := utils.Join(assets, "config/git/app.ini")

	if _, err = os.Stat(configoutput); errors.Is(err, os.ErrNotExist) {
		err = utils.Render(configtemplate, configoutput, gitconfigparams)
	}

	return err
}

func GetGitParams(assets string, runtype string) (*ProgramParams, error) {

	exefile := utils.Join(assets, "bin/git/gitea.exe")

	return &ProgramParams{
		ID:          "gitea",
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"web",
			"-c",
			utils.Join(assets, "config/git/app.ini"),
			"-w",
			utils.Join(assets, "data/git"),
			"--verbose",
		},
		LogFile: utils.Join(assets, "logs/git/git."+runtype+".log"),
	}, nil
}

func ExecP() {

}
