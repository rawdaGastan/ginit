package internal

import (
	"errors"
	"os"
	"strings"
)

// load files
func Readfile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// parse env file content
func Loadenv(content string) []string {
	envs := []string{}
	for _, line := range strings.Split(string(content), "\n") {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}
		key, value := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		envs = append(envs, key+"="+value)
	}
	return envs
}

// it parses the procfile and converts it to a list of procs
func LoadProcfile(content string) ([]*ProcInstance, error) {

	procs := []*ProcInstance{}
	for _, line := range strings.Split(string(content), "\n") {
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 || tokens[0][0] == '#' {
			continue
		}
		key, value := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])

		proc := &ProcInstance{name: key, cmdline: value}
		procs = append(procs, proc)
	}

	if len(procs) == 0 {
		return procs, errors.New("invalid entry, no procs found")
	}

	return procs, nil
}
