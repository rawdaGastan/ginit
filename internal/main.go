package internal

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type ProcInstance struct {
	name    string
	cmdline string
	cmd     *exec.Cmd
	index   int
}

type GinitService struct {
	procfile    string `yaml:"procfile"`
	mutex       sync.Mutex
	procs       []*ProcInstance
	args        []string
	exitOnError bool
}

// NewGinitService creates new instance from the GinitService
func NewGinitService(procfile string) GinitService {
	return GinitService{
		procfile:    procfile,
		exitOnError: true,
	}
}

func (gs *GinitService) ReadProcfile(procfile string) error {

	index := 0

	if procfile == "" {
		procfile = gs.procfile
	}

	content, err := os.ReadFile(procfile)
	if err != nil {
		return err
	}
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	for _, line := range strings.Split(string(content), "\n") {
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 || tokens[0][0] == '#' {
			continue
		}
		key, value := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])

		proc := &ProcInstance{name: key, cmdline: value, index: index}
		gs.procs = append(gs.procs, proc)
		index = index + 1
	}
	if len(gs.procs) == 0 {
		return errors.New("invalid entry")
	}

	return nil
}

// For start command
func (gs *GinitService) Start(ctx context.Context, sig <-chan os.Signal) error {
	err := gs.ReadProcfile("")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	// Cancel the RPC server when procs have returned/errored
	// cancel the context anyway in case of early return.
	defer cancel()

	if len(gs.args) > 1 {
		saveProcs := make([]*ProcInstance, 0, len(gs.args[1:]))
		for _, procName := range gs.args[1:] {
			proc := gs.getProc(procName)
			if proc == nil {
				return errors.New("proc is not found: " + procName)
			}
			saveProcs = append(saveProcs, proc)
		}
		gs.mutex.Lock()
		gs.procs = saveProcs
		gs.mutex.Unlock()
	}
	//godotenv.Load()

	//procsErr := gs.startProcs(sig, gs.exitOnError)
	//return procsErr
	return nil
}
