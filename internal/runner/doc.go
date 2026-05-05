// Package runner provides functionality to execute shell commands as cron job
// workloads and capture their results, including stdout, stderr, exit code,
// execution duration, and any errors encountered.
//
// Basic usage:
//
//	r := runner.New(30 * time.Second) // 30s timeout
//	result := r.Run("/usr/local/bin/my-job", "--flag", "value")
//	if !result.Success() {
//		log.Printf("job failed: exit=%d stderr=%s", result.ExitCode, result.Stderr)
//	}
//
// A zero timeout means the command runs without a deadline.
package runner

import "strings"
