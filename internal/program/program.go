package program

import (
	"os"
	"os/exec"
	"sync"

	"github.com/kardianos/service"
	"github.com/pnocera/novaco/internal/cmdparams"
	"github.com/pnocera/novaco/internal/utils"
)

var (
	writingSync sync.Mutex
	processes   []*exec.Cmd
)

type Program struct {
	Service service.Service
	Logger  service.Logger
}

func (p Program) Start(s service.Service) error {
	p.Logger.Info(s.String() + " starting")

	go p.run()
	return nil
}

func (p Program) Stop(s service.Service) error {
	p.Logger.Info(s.String() + " stopping")
	writingSync.Lock()
	for _, process := range processes {
		err := process.Process.Kill()
		if err != nil {
			p.Logger.Error(err)
		}
	}

	writingSync.Unlock()
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p Program) run() error {

	p.Logger.Info("Starting")
	var progparams []cmdparams.ProgramParams

	assets := utils.CurPath("assets")

	nomadparams, err := cmdparams.GetNomadProgramParams(assets)
	if err != nil {
		p.Logger.Error(err)
		return err
	}

	consulparams, err := cmdparams.GetConsulProgramParams(assets)
	if err != nil {
		p.Logger.Error(err)
		return err
	}

	progparams = append(progparams, *consulparams, *nomadparams)

	return p.ExecAndWait(progparams)

}

func (p Program) ExecAndWait(commands []cmdparams.ProgramParams) error {
	var wg sync.WaitGroup

	for _, param := range commands {
		cmd := exec.Command(param.ExeFullname, param.AdditionalParams...)
		cmd.Dir = param.DirPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		writingSync.Lock()
		processes = append(processes, cmd)
		writingSync.Unlock()

		wg.Add(1)

		go func() {
			defer wg.Done()
			err := cmd.Start()
			if err != nil {
				p.Logger.Error(err)
			}

		}()
	}
	wg.Wait()

	return nil
}
