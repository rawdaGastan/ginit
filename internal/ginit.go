package internal

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type ProcInstance struct {
	name    string
	cmdline string
	cmd     *exec.Cmd
	stopped bool
}

type GinitService struct {
	procfile string `yaml:"procfile"`
	envfile  string
	procs    []*ProcInstance

	// args are the specidied procs coming after start cmd, for example: start web
	args []string

	exitOnError bool
	sig         chan os.Signal
	logger      zerolog.Logger
	env         []string
}

// NewGinitService creates new instance from the GinitService
func NewGinitService(procfile string, envfile string, args []string, logger zerolog.Logger) GinitService {
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	return GinitService{
		procfile:    procfile,
		envfile:     envfile,
		args:        args,
		exitOnError: true,
		sig:         sig,
		logger:      logger,
	}
}

// For start command
func (gs *GinitService) Start(ctx context.Context) error {
	// read procfile
	content, err := Readfile(gs.procfile)
	if err != nil {
		return err
	}

	gs.procs, err = LoadProcfile(content)
	if err != nil {
		return err
	}

	// read envfile
	content, err = Readfile(gs.envfile)
	if err != nil {
		return err
	}
	gs.env = Loadenv(content)

	// start
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if len(gs.args) > 0 {
		saveProcs := make([]*ProcInstance, 0, len(gs.args))
		for _, procName := range gs.args {
			proc := gs.getProc(procName)
			if proc == nil {
				return errors.New("proc is not found: " + procName)
			}
			saveProcs = append(saveProcs, proc)
		}
		gs.procs = saveProcs
	}

	procsErr := gs.startProcs()
	return procsErr
}
