package newprogram

import (
	"fmt"
	"log"
	"os"

	cmd "github.com/ShinyTrinkets/overseer"

	"github.com/go-ini/ini"
	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/cmdparams"
	"github.com/pnocera/novaco/internal/utils"
)

type newprogram struct {
	runtype string
	ovr     *cmd.Overseer
}

func NewProgram(runtype string) *newprogram {

	return &newprogram{runtype: runtype, ovr: cmd.NewOverseer()}
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

	assets, err := utils.Assets()
	if err != nil {
		return err
	}

	err = cmdparams.RenderGitConfigIfNotExist(p.runtype)
	if err != nil {
		return err
	}

	gitparams, err := cmdparams.GetGitParams(assets, p.runtype)
	if err != nil {
		return err
	}

	nomadparams, err := cmdparams.GetNomadProgramParams(assets, p.runtype)
	if err != nil {
		return err
	}

	consulparams, err := cmdparams.GetConsulProgramParams(assets, p.runtype)
	if err != nil {
		return err
	}

	vaultparams, err := cmdparams.GetVaultProgramParams(assets, p.runtype)
	if err != nil {
		return err
	}

	progparams = append(
		progparams,
		// *consulparams,
		// *vaultparams,
		// *nomadparams,
		*gitparams)

	log.Printf("consulparams: %v", *consulparams)
	log.Printf("vaultparams: %v", *vaultparams)
	log.Printf("nomadparams: %v", *nomadparams)
	log.Printf("gitparams: %v", *gitparams)

	err = p.ExecAndWait(progparams)

	return err

}

// type DummyLogger struct {
// 	Name string
// }

// func (l *DummyLogger) Info(msg string, v ...interface{}) {
// 	log.Println(msg, v)
// }
// func (l *DummyLogger) Error(msg string, v ...interface{}) {
// 	log.Println(msg, v)
// }

func (p newprogram) ExecAndWait(commands []cmdparams.ProgramParams) error {
	cmd.SetupLogBuilder(func(name string) cmd.Logger {
		logger := NewKLogger(name)
		return logger
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

	go func() {
		for state := range statusFeed {

			if state.ID == "gitea" {
				if state.State == "running" {
					// wait for health check
					gitini, err := cmdparams.GetGitIniFilePath(p.runtype)
					if err != nil {
						log.Printf("%v", err)
					} else {
						cfg, err := ini.Load(gitini)
						if err != nil {
							log.Printf("%v", err)
						} else {
							url := cfg.Section("server").Key("ROOT_URL").String()
							err = utils.WaitForUrl(url)
							if err != nil {
								log.Printf("%v", err)
							} else {
								err = utils.InitGitea(p.runtype)
								if err != nil {
									log.Printf("%v", err)
								}
							}
						}
					}

				}
			}

		}
	}()

	logFeed := make(chan *cmd.LogMsg)
	p.ovr.WatchLogs(logFeed)

	go func() {
		for logmsg := range logFeed {
			fmt.Printf("LOG: %v\n", logmsg)
		}
	}()

	p.ovr.SuperviseAll()

	return nil
}

func NewKLogger(name string) {
	panic("unimplemented")
}
