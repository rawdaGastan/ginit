package internal

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
)

func BenchmarkTestStartGinit(b *testing.B) {

	b.Run("test_start", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger := zerolog.New(os.Stdout).With().Logger()
			procfile := "../demo/Procfile"
			envfile := "../demo/.env"

			ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
			err := ginit.Start(context.Background())

			if err != nil {
				b.Errorf("start ginit service should have no errors")
			}
		}
	})
}

func TestHelpers(t *testing.T) {

	t.Run("readfile", func(t *testing.T) {
		_, err := Readfile("../demo/.env")

		if err != nil {
			t.Errorf("file should generate no errors, %v", err)
		}
	})

	t.Run("read_wrong_file", func(t *testing.T) {
		_, err := Readfile("demo/.env")

		if err == nil {
			t.Errorf("file should generate an error")
		}
	})

	t.Run("read_content_from_file", func(t *testing.T) {
		content, _ := Readfile("../demo/.env")

		want := "USERNAME=rawda\nPASSWORD=elpassword\nREDIS_PORT=6379\nPY_PORT=8080"

		if content != want {
			t.Errorf("content should be, %v", want)
		}
	})

	t.Run("load_envs", func(t *testing.T) {
		content := "USERNAME=rawda\nPASSWORD=elpassword"
		envs := Loadenv(content)

		if len(envs) != 2 {
			t.Errorf("envs count should be 2")
		}
	})

	t.Run("load_wrong_envs", func(t *testing.T) {
		content := "USERNAME=user=rawda\nPASSWORD=elpassword"
		envs := Loadenv(content)

		if len(envs) != 1 {
			t.Errorf("envs count should be 1")
		}
	})

	t.Run("load_procs", func(t *testing.T) {
		content := "redis: redis-server --port $REDIS_PORT"
		procs, _ := LoadProcfile(content)

		want := &ProcInstance{name: "redis", cmdline: "redis-server --port $REDIS_PORT"}

		if procs[0].name != want.name {
			t.Errorf("proc name should be redis")
		}
	})

	t.Run("load_comment_procs", func(t *testing.T) {
		content := "#"
		procs, _ := LoadProcfile(content)

		if len(procs) != 0 {
			t.Errorf("procs count should be 0")
		}
	})

	t.Run("load_wrong_procs", func(t *testing.T) {
		content := ""
		_, err := LoadProcfile(content)

		if err == nil {
			t.Errorf("procs should generate an error")
		}
	})
}

func TestService(t *testing.T) {

	t.Run("test_new_service", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{}, logger)

		if ginit.args == nil {
			t.Errorf("ginit service should have empty args")
		}
	})

	t.Run("test_start", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
		err := ginit.Start(context.Background())

		if err != nil {
			t.Errorf("start ginit service should have no errors")
		}
	})

	t.Run("test_start_wrong_proc", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test2"}, logger)
		err := ginit.Start(context.Background())

		if err == nil {
			t.Errorf("start ginit service should have an error, test2 is not found")
		}
	})

	t.Run("test_start_wrong_env", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "/demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
		err := ginit.Start(context.Background())

		if err == nil {
			t.Errorf("start ginit service should have an error, envfile is not found")
		}
	})

	t.Run("test_start_wrong_proc", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "/demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
		err := ginit.Start(context.Background())

		if err == nil {
			t.Errorf("start ginit service should have an error, procfile is not found")
		}
	})

	t.Run("test_start_wrong_proc_cmd", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"wrong"}, logger)
		procs := []*ProcInstance{{name: "wrong", cmdline: "web"}}
		ginit.procs = procs
		err := ginit.startProcs()

		if err == nil {
			t.Errorf("start procs in ginit service should have an error")
		}
	})

	t.Run("test_start_interrupt_sig", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
		ginit.sig <- os.Interrupt
		err := ginit.startProcs()

		if err != nil {
			t.Errorf("start procs in ginit service should have no errors")
		}
	})

	t.Run("test_stop_interrupt_sig", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"web"}, logger)
		procs := []*ProcInstance{{name: "web", cmdline: "python3 -m http.server $PY_PORT"}}
		ginit.procs = procs
		err := ginit.stopProc("web", nil)

		if err != nil {
			t.Errorf("stop proc in ginit service should have no errors %v", err)
		}
	})

	t.Run("test_stop_proc_kill", func(t *testing.T) {
		logger := zerolog.New(os.Stdout).With().Logger()
		procfile := "../demo/Procfile"
		envfile := "../demo/.env"

		ginit := NewGinitService(procfile, envfile, []string{"test"}, logger)
		cmd := exec.Command("/bin/sh", "-c echo")
		err := cmd.Start()

		if err != nil {
			t.Errorf("an error occurred %v", err)
		}

		procs := []*ProcInstance{{name: "test", cmdline: "echo", cmd: cmd}}
		ginit.procs = procs

		err = ginit.stopProc("test", os.Kill)

		if err != nil {
			t.Errorf("stop proc in ginit service should have no errors %v", err)
		}
	})
}
