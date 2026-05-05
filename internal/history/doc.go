// Package history implements run-history tracking for cronwrap.
//
// It provides a thread-safe, file-backed store that records the outcome of
// each cron job execution, including timing, exit code, stdout, and stderr.
//
// Usage:
//
//	store, err := history.NewStore("/var/lib/cronwrap/history.json")
//	if err != nil { ... }
//
//	err = store.Append(history.Record{
//		JobName:  "my-job",
//		Command:  "./backup.sh",
//		Success:  true,
//		ExitCode: 0,
//	})
//
// Records are stored as a JSON array and are safe for concurrent use.
package history
