package program

import (
	"log"
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
	errs        chan error
)

type program struct {
	runtype string
	Logger  service.Logger
}

func NewProgram(runtype string) *program {

	return &program{runtype: runtype}
}

func (p program) Start(s service.Service) error {
	logger, err := s.Logger(nil)

	if err != nil {
		log.Fatal(err)
	}
	p.Logger = logger
	p.Logger.Info(s.String() + " starting")
	errs = make(chan error, 1)
	go p.run()
	return <-errs
}

func (p program) Stop(s service.Service) error {
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

func (p program) run() {

	p.Logger.Info("Starting")
	var progparams []cmdparams.ProgramParams

	assets, err := utils.Assets()
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	err = cmdparams.RenderGitConfigIfNotExist(p.runtype)
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	gitparams, err := cmdparams.GetGitParams(assets, p.runtype)
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	nomadparams, err := cmdparams.GetNomadProgramParams(assets, p.runtype)
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	consulparams, err := cmdparams.GetConsulProgramParams(assets, p.runtype)
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	vaultparams, err := cmdparams.GetVaultProgramParams(assets, p.runtype)
	if err != nil {
		p.Logger.Error(err)
		errs <- err
		close(errs)
		return
	}

	progparams = append(progparams, *consulparams, *vaultparams, *nomadparams, *gitparams)
	p.Logger.Info(consulparams.ExeFullname)
	p.Logger.Info(vaultparams.ExeFullname)
	p.Logger.Info(nomadparams.ExeFullname)
	p.Logger.Info(gitparams.ExeFullname)

	err = p.ExecAndWait(progparams)
	if err != nil {
		p.Logger.Error(err)
		errs <- err

	}

	close(errs)

}

func (p program) ExecAndWait(commands []cmdparams.ProgramParams) error {
	var wg sync.WaitGroup
	var suberrs chan error = make(chan error, len(commands))

	for _, param := range commands {
		cmd := exec.Command(param.ExeFullname, param.AdditionalParams...)
		cmd.Dir = param.DirPath

		var f *os.File = nil
		if param.LogFile != "" {
			f, _ = os.OpenFile(param.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		}
		// if f != nil {
		// 	defer f.Close()
		// }

		cmderr := os.Stderr
		if f != nil {
			cmderr = f
		}

		cmdout := os.Stdout
		if f != nil {
			cmdout = f
		}

		cmd.Stdout = cmdout
		cmd.Stderr = cmderr

		writingSync.Lock()
		processes = append(processes, cmd)
		writingSync.Unlock()

		wg.Add(1)

		go func() {
			defer wg.Done()
			err0 := cmd.Start()
			if err0 != nil {
				p.Logger.Error(err0)
				suberrs <- err0
			}

		}()
	}

	if <-suberrs != nil {
		return <-suberrs
	}

	wg.Wait()

	return nil
}
