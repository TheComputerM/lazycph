package core

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrCompile = errors.New("compilation error")
	ErrExecute = errors.New("execution error")
)

type Engine struct {
	// mode can be 'compile' or 'interpret'
	Mode    string   `json:"mode"`
	Command []string `json:"command"`
}

func (e *Engine) compile(path, input string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "lazycph-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "program")

	// Replace placeholders in each argument
	replacer := strings.NewReplacer("{file}", path, "{temp}", outputPath)
	args := make([]string, len(e.Command))
	for i, arg := range e.Command {
		args[i] = replacer.Replace(arg)
	}

	// Compile
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return string(out), ErrCompile
	}

	// Run
	runCmd := exec.Command(outputPath)
	runCmd.Stdin = strings.NewReader(input)
	out, err = runCmd.CombinedOutput()
	if err != nil {
		return string(out), ErrExecute
	}

	return string(out), nil
}

func (e *Engine) interpret(path, input string) (string, error) {
	replacer := strings.NewReplacer("{file}", path)
	args := make([]string, len(e.Command))
	for i, arg := range e.Command {
		args[i] = replacer.Replace(arg)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(input)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), ErrExecute
	}

	return string(out), nil
}

func (e *Engine) Run(path, input string) (string, error) {
	if e.Mode == "compile" {
		return e.compile(path, input)
	}

	return e.interpret(path, input)
}
