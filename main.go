package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var testProjectPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/seriousben/is-it-binary/")

func runDep(projectDir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(context.TODO(), "dep", args...)
	cmd.Dir = projectDir

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return cmd.Output()
}

func requiresUpdate() ([]byte, error) {
	data, err := runDep(testProjectPath, strings.Split("ensure -update --dry-run", " ")...)
	if err == nil {
		if len(data) == 0 {
			return data, errors.New("does not require an update")
		}
		return data, err
	}
	return data, err
}

func update() ([]byte, error) {
	return runDep(testProjectPath, strings.Split("ensure -update", " ")...)
}

func printVersion() ([]byte, error) {
	return runDep(".", "version")
}

func printError() ([]byte, error) {
	return runDep(".", "notacommand")
}

func do(label string, cmd func() ([]byte, error)) {
	fmt.Printf("=> Start %s\n", label)
	defer fmt.Printf("<= End %s\n", label)
	data, err := cmd()
	fmt.Printf("data: %s\n", string(data))
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Printf("exit-error: %+v\n", ee, string(ee.Stderr))
		} else {
			fmt.Printf("error: %+v\n", err)
		}
	}
}

func main() {
	do("printError", printError)
	do("printVersion", printVersion)
	do("requiresUpdate", requiresUpdate)
}
