package runner

import (
	"context"
	"os/exec"
	"time"
)

// Result holds the outcome of a cron job execution.
type Result struct {
	Command   string
	Args      []string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	ExitCode  int
	Stdout    string
	Stderr    string
	Error     error
}

// Runner executes shell commands and captures their output.
type Runner struct {
	Timeout time.Duration
}

// New creates a new Runner with the given timeout.
// If timeout is zero, no timeout is applied.
func New(timeout time.Duration) *Runner {
	return &Runner{Timeout: timeout}
}

// Run executes the given command with arguments and returns a Result.
func (r *Runner) Run(command string, args ...string) *Result {
	result := &Result{
		Command:   command,
		Args:      args,
		StartTime: time.Now(),
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if r.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()
	result.Error = err

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	} else if err != nil {
		result.ExitCode = 1
	}

	return result
}

// Success returns true if the command exited with code 0.
func (r *Result) Success() bool {
	return r.ExitCode == 0
}
