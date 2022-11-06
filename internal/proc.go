package internal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

// getProc gets the proc intance using its name
func (gs *GinitService) getProc(name string) *ProcInstance {
	for _, proc := range gs.procs {
		if proc.name == name {
			return proc
		}
	}
	return nil
}

// start specified proc. if proc is started already, return nil.
func (gs *GinitService) startProc(name string, wg *sync.WaitGroup, errCh chan<- error) error {
	proc := gs.getProc(name)
	if proc == nil {
		return errors.New("unknown name: " + name)
	}

	if proc.cmd != nil {
		return nil
	}

	if wg != nil {
		wg.Add(1)
	}
	go func() {
		gs.spawnProc(name, errCh)
		if wg != nil {
			wg.Done()
		}
	}()
	return nil
}

// spawnProc starts the specified proc, and returns any error from running it.
func (gs *GinitService) spawnProc(name string, errCh chan<- error) {
	proc := gs.getProc(name)

	gs.logger.Info().Msg("\nexecuting " + proc.cmdline + "\n")

	cmdToExecute := []string{"-c", proc.cmdline}

	cmd := exec.Command("/bin/sh", cmdToExecute...)

	cmd.Env = gs.env
	cmd.Stdin = nil
	cmd.Stdout = gs.logger
	cmd.Stderr = gs.logger

	if err := cmd.Start(); err != nil {
		select {
		case errCh <- err:
		default:
		}
		gs.logger.Error().Msg("failed to start " + name + ": " + fmt.Sprint(err) + "\n")
		return
	}

	proc.cmd = cmd
	proc.stopped = false

	err := cmd.Wait()
	if err != nil && !proc.stopped {
		select {
		case errCh <- err:
		default:
		}
	}

	// proc cmd is done
	proc.cmd = nil
	gs.logger.Info().Msg("terminating " + name + "\n")
}

// Stop the specified proc, If signal is nil, os.Interrupt is used.
func (gs *GinitService) stopProc(name string, signal os.Signal) error {
	if signal == nil {
		signal = os.Interrupt
	}

	proc := gs.getProc(name)
	if proc == nil {
		return errors.New("not found proc: " + name)
	}

	if proc.cmd == nil {
		return nil
	}
	proc.stopped = true

	err := proc.cmd.Process.Signal(signal)
	if err != nil {
		return err
	}

	timeout := time.AfterFunc(10*time.Second, func() {
		if proc.cmd != nil {
			err = proc.cmd.Process.Kill()
		}
	})
	timeout.Stop()
	return err
}

// stopProcs attempts to stop every running process,
// stopProcs will wait until all procs have had an opportunity to stop.
func (gs *GinitService) stopProcs(sig os.Signal) error {
	var err error
	for _, proc := range gs.procs {
		stopErr := gs.stopProc(proc.name, sig)
		if stopErr != nil {
			err = stopErr
		}
	}
	return err
}

// spawn all procs
func (gs *GinitService) startProcs() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, proc := range gs.procs {
		gs.startProc(proc.name, &wg, errCh)
	}

	allProcsDone := make(chan bool, 1)
	go func() {
		wg.Wait()
		allProcsDone <- true
	}()

	for {
		select {
		case err := <-errCh:
			if gs.exitOnError {
				gs.stopProcs(os.Interrupt)
				return err
			}
		case <-allProcsDone:
			return gs.stopProcs(os.Interrupt)
		case sig := <-gs.sig:
			return gs.stopProcs(sig)
		}
	}
}
