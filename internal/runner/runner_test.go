package runner_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cronwrap/cronwrap/internal/runner"
)

func TestRunSuccess(t *testing.T) {
	r := runner.New(5 * time.Second)
	result := r.Run("echo", "hello")

	if !result.Success() {
		t.Errorf("expected success, got exit code %d", result.ExitCode)
	}
	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("expected stdout to contain 'hello', got %q", result.Stdout)
	}
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestRunFailure(t *testing.T) {
	r := runner.New(5 * time.Second)
	result := r.Run("false")

	if result.Success() {
		t.Error("expected failure, got success")
	}
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit code")
	}
}

func TestRunTimeout(t *testing.T) {
	r := runner.New(100 * time.Millisecond)
	result := r.Run("sleep", "5")

	if result.Success() {
		t.Error("expected timeout failure")
	}
	if result.Error == nil {
		t.Error("expected an error due to timeout")
	}
}

func TestRunStderr(t *testing.T) {
	r := runner.New(5 * time.Second)
	result := r.Run("sh", "-c", "echo error >&2; exit 1")

	if result.Success() {
		t.Error("expected failure")
	}
	if !strings.Contains(result.Stderr, "error") {
		t.Errorf("expected stderr to contain 'error', got %q", result.Stderr)
	}
}

func TestRunTimestamps(t *testing.T) {
	r := runner.New(0)
	before := time.Now()
	result := r.Run("echo", "ts")
	after := time.Now()

	if result.StartTime.Before(before) || result.StartTime.After(after) {
		t.Error("StartTime out of expected range")
	}
	if result.EndTime.Before(result.StartTime) {
		t.Error("EndTime before StartTime")
	}
}
