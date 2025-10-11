package exec

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"forge/src/logging"
	"forge/src/output"
	"os/exec"
	"strings"
	"time"
)

// ExecutionStep minimal struct for hanging detection test.
type ExecutionStep struct {
	Start      time.Time
	LastOutput time.Time
	Cmd        PrimitiveCommand
	ExitCode   int
}

// IsHanging returns true if >60s since LastOutput.
func (s ExecutionStep) IsHanging(now time.Time) bool {
	if s.LastOutput.IsZero() {
		return false
	}
	return now.Sub(s.LastOutput) > 60*time.Second
}

// ValidateCommand returns error for non-whitelist commands.
func ValidateCommand(cmd string) error {
	if !IsWhitelisted(cmd) {
		return errors.New("ERR_NON_WHITELIST: " + cmd)
	}
	return nil
}

// RunSingleCommand executes a whitelisted command and streams output lines into buffer.
func RunSingleCommand(ctx context.Context, primitive PrimitiveCommand, buf *output.OutputBuffer) (exitCode int, err error) {
	if err = ValidateCommand(primitive.Cmd); err != nil {
		return -1, err
	}
	logging.LogEvent(logging.EventCommandStart, "cmd", primitive.Cmd)
	parts := strings.Split(primitive.Cmd, " ")
	name := parts[0]
	args := parts[1:]
	c := exec.CommandContext(ctx, name, args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return -1, err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return -1, err
	}
	if err = c.Start(); err != nil {
		return -1, err
	}
	// stream output concurrently
	go func() {
		s := bufio.NewScanner(stdout)
		for s.Scan() {
			buf.Append(s.Text())
		}
	}()
	go func() {
		s := bufio.NewScanner(stderr)
		for s.Scan() {
			buf.Append(s.Text())
		}
	}()
	err = c.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			logging.LogEvent(logging.EventCommandEnd, "cmd", primitive.Cmd, "exit", exitErr.ExitCode())
			return exitErr.ExitCode(), nil
		}
		return -1, err
	}
	logging.LogEvent(logging.EventCommandEnd, "cmd", primitive.Cmd, "exit", 0)
	return 0, nil
}

// Workflow represents a sequential set of primitive commands.
type Workflow struct {
	Steps []ExecutionStep
}

// RunWorkflow executes steps sequentially; aborts on first non-zero exit code.
func RunWorkflow(ctx context.Context, wf *Workflow, buf *output.OutputBuffer) (failedAt int, err error) {
	logging.LogEvent(logging.EventWorkflowStart, "steps", len(wf.Steps))
	if len(wf.Steps) == 0 {
		// structural validation failure (T084): zero-step workflows are rejected early
		logging.LogEvent(logging.EventWorkflowEnd, "failedIndex", -1, "reason", "empty")
		return -1, ErrEmptyWorkflow
	}
	for i := range wf.Steps {
		step := &wf.Steps[i]
		step.Start = time.Now()
		exit, runErr := RunSingleCommand(ctx, step.Cmd, buf)
		step.ExitCode = exit
		step.LastOutput = time.Now()
		if runErr != nil {
			buf.Append("STEP FAILED (non-whitelist): " + step.Cmd.Cmd + " color=" + string(output.ErrorColor))
			logging.LogEvent(logging.EventWorkflowEnd, "failedIndex", i)
			return i, runErr
		}
		if exit != 0 {
			buf.Append("STEP FAILED (exit code) " + step.Cmd.Cmd + " code=" + fmtCode(exit) + " color=" + string(output.ErrorColor))
			logging.LogEvent(logging.EventWorkflowEnd, "failedIndex", i)
			return i, errors.New("workflow step failed")
		}
	}
	logging.LogEvent(logging.EventWorkflowEnd, "failedIndex", -1)
	return -1, nil
}

func fmtCode(c int) string { return fmt.Sprintf("%d", c) }

// ErrEmptyWorkflow indicates a workflow with zero steps is not executable.
var ErrEmptyWorkflow = errors.New("workflow has no steps")
