package terminal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RunCmdOptions struct {
	Cwd        string
	Env        []string
	ShowOutput bool
	CtxTimeout time.Duration
}

// RunCommand runs the given script, streaming stdout and stderr to the
// terminal and capturing exit error results in the returned results.
func RunCommand(script string, opt *RunCmdOptions) ([]byte, error) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if opt.CtxTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opt.CtxTimeout)
		defer cancel()
	}

	if containsSudo(script) {
		return RunCmdAndExpectUserInput(script, opt)
	}

	if opt.ShowOutput {
		fmt.Printf("running script: %s\n", script)
	}

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)

	// make sure the script runs with the current environment
	// this allows things like PATH setting to work accoss shells
	cmd.Env = os.Environ()

	if len(opt.Env) > 0 {
		cmd.Env = append(cmd.Env, opt.Env...)
	}

	if opt.Cwd != "" {
		cmd.Dir = opt.Cwd
	}

	if opt.ShowOutput {
		results := bytes.Buffer{}
		cmd.Stdout = io.MultiWriter(&results, os.Stdout)

		detailedErr := bytes.Buffer{}
		cmd.Stderr = io.MultiWriter(os.Stderr, &detailedErr)

		err := cmd.Run()
		if err != nil {
			return results.Bytes(), doErr(err, script, detailedErr.String())
		}
		return results.Bytes(), nil
	}
	results := bytes.Buffer{}
	cmd.Stdout = &results
	detailedErr := bytes.Buffer{}
	cmd.Stderr = io.MultiWriter(&detailedErr)

	err := cmd.Run()
	if err != nil {
		return nil, doErr(err, script, detailedErr.String())
	}

	return results.Bytes(), nil
}

func RunCmdAndExpectUserInput(script string, opt *RunCmdOptions) ([]byte, error) {
	if opt.CtxTimeout <= 0 {
		opt.CtxTimeout = time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), opt.CtxTimeout)
	defer cancel()

	return RunCmdContextAndExpectUserInput(ctx, script, opt)
}

func RunCmdContextAndExpectUserInput(ctx context.Context, script string, opt *RunCmdOptions) ([]byte, error) {
	if opt.ShowOutput {
		fmt.Printf("Running: %s\n", script)
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)
	cmd.Stdin = os.Stdin // This will cause the command to pause if there is a prompt waiting for stdin

	cmd.Env = os.Environ()

	if len(opt.Env) > 0 {
		cmd.Env = append(cmd.Env, opt.Env...)
	}

	if opt.Cwd != "" {
		cmd.Dir = opt.Cwd
	}

	if opt.ShowOutput {
		cmd.Stdout = os.Stdout

		detailedErr := bytes.Buffer{}
		cmd.Stderr = io.MultiWriter(os.Stderr, &detailedErr)

		err := cmd.Run()
		if err != nil {
			return nil, doErr(err, script, detailedErr.String())
		}
		return nil, nil
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, doErr(err, script, string(out))
	}

	return out, nil
}

func containsSudo(script string) bool {
	return strings.Contains(script, "sudo")
}

func isSignalInterrupt(err error) bool {
	return strings.Contains(err.Error(), "signal: interrupt")
}

func doErr(err error, script, message string) error {
	if err != nil && isSignalInterrupt(err) {
		return fmt.Errorf("failed to run script -> %s | error code: %s | error message: %s", script, err.Error(), message)
	} else if err != nil && !errors.Is(err, io.ErrShortWrite) {
		return fmt.Errorf("failed to run script -> %s | error code: %s | error message: %s", script, err.Error(), message)
	}
	return nil
}
