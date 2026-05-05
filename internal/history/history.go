// Package history provides run-history tracking for cronwrap.
package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Record represents a single cron job execution record.
type Record struct {
	JobName   string        `json:"job_name"`
	Command   string        `json:"command"`
	StartedAt time.Time     `json:"started_at"`
	FinishedAt time.Time    `json:"finished_at"`
	Duration  time.Duration `json:"duration_ns"`
	ExitCode  int           `json:"exit_code"`
	Success   bool          `json:"success"`
	Stdout    string        `json:"stdout,omitempty"`
	Stderr    string        `json:"stderr,omitempty"`
}

// Store persists run history records to a JSON file.
type Store struct {
	mu      sync.Mutex
	path    string
	records []Record
}

// NewStore creates a new Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Append adds a new record and persists the store to disk.
func (s *Store) Append(r Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, r)
	return s.save()
}

// Records returns a copy of all stored records.
func (s *Store) Records() []Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Record, len(s.records))
	copy(out, s.records)
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.records)
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
