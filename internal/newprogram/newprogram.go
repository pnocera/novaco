package newprogram

import (
	"os"
	"strings"

	cmd "github.com/ShinyTrinkets/overseer"

	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/cmdparams"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/utils"
)

var (
	sets = settings.GetSettings()
)

type newprogram struct {
	ovr *cmd.Overseer
}

func NewProgram() *newprogram {

	return &newprogram{ovr: cmd.NewOverseer()}
}

func (p newprogram) Start(s service.Service) error {
	return p.run()

}

func (p newprogram) Stop(s service.Service) error {

	p.ovr.StopAll(false)
	return nil
}

func (p newprogram) run() error {

	var progparams []cmdparams.ProgramParams

	err := cmdparams.RenderGitConfigIfNotExist()
	if err != nil {
		sets.Logger.Error("Error rendering git config %v", err)
		return err
	}

	gitparams, err := cmdparams.GetGitParams()
	if err != nil {
		sets.Logger.Error("Error getting git params %v", err)
		return err
	}

	csiproxyparams, err := cmdparams.GetProxyParams()
	if err != nil {
		sets.Logger.Error("Error getting csi-proxy params %v", err)
		return err
	}

	nomadparams, err := cmdparams.GetNomadProgramParams()
	if err != nil {
		sets.Logger.Error("Error getting nomad params %v", err)
		return err
	}

	consulparams, err := cmdparams.GetConsulProgramParams()
	if err != nil {
		sets.Logger.Error("Error getting consul params %v", err)
		return err
	}

	vaultparams, err := cmdparams.GetVaultProgramParams()
	if err != nil {
		sets.Logger.Error("Error getting vault params %v", err)
		return err
	}
	runtypes := strings.Split(sets.Runtypes, ",")
	if utils.StringsContains(runtypes, "server") {
		progparams = append(progparams, *gitparams, *csiproxyparams, *nomadparams, *consulparams, *vaultparams)
	}

	err = p.ExecAndWait(progparams)
	if err != nil {
		sets.Logger.Error("Error executing programs %v", err)
	}

	return err

}

func (p newprogram) ExecAndWait(commands []cmdparams.ProgramParams) error {

	cmd.SetupLogBuilder(func(name string) cmd.Logger {
		return settings.NewKLogger(name)

	})

	opts := cmd.Options{
		Buffered: false, Streaming: true,
		Env: os.Environ(),
	}

	for _, command := range commands {
		p.ovr.Add(command.ID, command.ExeFullname, command.AdditionalParams, opts)
	}

	statusFeed := make(chan *cmd.ProcessJSON)
	p.ovr.WatchState(statusFeed)

	utils.WatchStatus(statusFeed)

	p.ovr.SuperviseAll()

	return nil
}
